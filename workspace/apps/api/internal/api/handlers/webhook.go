package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/emiliospot/footie/api/internal/config"
	"github.com/emiliospot/footie/api/internal/infrastructure/events"
	"github.com/emiliospot/footie/api/internal/infrastructure/webhooks"
	"github.com/emiliospot/footie/api/internal/repository/sqlc"
)

// WebhookHandler handles webhook endpoints for external event providers.
type WebhookHandler struct {
	*BaseHandler
	webhookConfig *config.WebhookConfig
	providerRegistry *webhooks.Registry
}

// NewWebhookHandler creates a new webhook handler.
func NewWebhookHandler(base *BaseHandler, webhookConfig *config.WebhookConfig, providerRegistry *webhooks.Registry) *WebhookHandler {
	return &WebhookHandler{
		BaseHandler:     base,
		webhookConfig:   webhookConfig,
		providerRegistry: providerRegistry,
	}
}

// ExternalEventPayload represents the incoming webhook payload from external providers.
type ExternalEventPayload struct {
	MatchID           int32   `json:"matchId" binding:"required"`
	EventType         string  `json:"eventType" binding:"required"` // GOAL, SHOT, PASS, CARD, SUBSTITUTION, etc.
	Minute            int32   `json:"minute" binding:"required"`
	ExtraMinute       *int32  `json:"extraMinute,omitempty"`
	TeamID            *int32  `json:"teamId,omitempty"`
	PlayerID          *int32  `json:"playerId,omitempty"`
	SecondaryPlayerID *int32  `json:"secondaryPlayerId,omitempty"`
	PositionX         *float64 `json:"positionX,omitempty"`
	PositionY         *float64 `json:"positionY,omitempty"`
	Description       string  `json:"description,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"` // xG, pass completion, etc.
	Timestamp         *string `json:"timestamp,omitempty"` // ISO 8601 format
}

// HandleMatchEvents handles POST /webhooks/matches.
// This endpoint receives match events from external providers (e.g., data feed services).
// Supports multiple providers via query parameter: ?provider=opta|statsbomb|generic
// @Summary Receive match events via webhook
// @Description Receives match events from external providers and processes them
// @Tags webhooks
// @Accept json
// @Produce json
// @Param provider query string false "Provider name (opta, statsbomb, generic)" default(generic)
// @Param X-Signature header string true "HMAC SHA256 signature"
// @Param X-Provider header string false "Provider identifier (alternative to query param)"
// @Param payload body map[string]interface{} true "Match event payload (provider-specific format)"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /webhooks/matches [post]
func (h *WebhookHandler) HandleMatchEvents(c *gin.Context) {
	// 1. Determine provider (from query param or header)
	providerName := c.Query("provider")
	if providerName == "" {
		providerName = c.GetHeader("X-Provider")
	}
	if providerName == "" {
		providerName = "generic" // Default to generic provider
	}
	providerName = strings.ToLower(providerName)

	// 2. Get provider from registry
	provider, err := h.providerRegistry.GetProvider(providerName)
	if err != nil {
		h.logger.Warn("Unknown provider", "provider", providerName, "ip", c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown provider", "available": h.providerRegistry.ListProviders()})
		return
	}

	// 3. Read raw payload for signature verification and extraction
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Warn("Failed to read request body", "error", err, "ip", c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read payload"})
		return
	}
	// Restore body for potential re-reading
	c.Request.Body = io.NopCloser(strings.NewReader(string(body)))

	// 4. Verify signature using provider-specific secret
	signature := c.GetHeader("X-Signature")
	secret := h.getProviderSecret(providerName)
	if !provider.VerifySignature(body, signature, secret) {
		h.logger.Warn("Invalid webhook signature", "provider", providerName, "ip", c.ClientIP())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	// 5. Extract events using provider-specific adapter (supports both single and batch)
	events, err := provider.ExtractEvents(c.Request.Context(), body)
	if err != nil {
		h.logger.Warn("Failed to extract events", "error", err, "provider", providerName, "ip", c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload format", "details": err.Error()})
		return
	}

	if len(events) == 0 {
		h.logger.Warn("No events extracted from payload", "provider", providerName, "ip", c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "No events found in payload"})
		return
	}

	// 6. Validate all matches exist and collect unique match IDs
	matchIDs := make(map[int32]bool)
	for _, event := range events {
		matchIDs[event.MatchID] = true
	}

	// Validate all matches exist
	for matchID := range matchIDs {
		_, err := h.queries.GetMatchByID(c.Request.Context(), matchID)
		if err != nil {
			h.logger.Error("Match not found for webhook event", "error", err, "match_id", matchID, "provider", providerName)
			c.JSON(http.StatusNotFound, gin.H{"error": "Match not found", "match_id": matchID})
			return
		}
	}

	// 7. Process all events asynchronously (store in DB + publish to Redis)
	go h.processProviderEventsAsync(c.Request.Context(), events, providerName)

	// 8. Acknowledge quickly (webhook best practice)
	response := gin.H{
		"status":      "accepted",
		"events_count": len(events),
		"provider":     providerName,
	}

	// Include match_id and event_type for single events (backward compatibility)
	if len(events) == 1 {
		response["match_id"] = events[0].MatchID
		response["event_type"] = events[0].EventType
	} else {
		// For batches, include summary
		eventTypes := make(map[string]int)
		for _, event := range events {
			eventTypes[event.EventType]++
		}
		response["event_types"] = eventTypes
	}

	c.JSON(http.StatusOK, response)
}

// processWebhookEventAsync processes the webhook event asynchronously.
// This includes storing in DB and publishing to Redis Streams/Pub/Sub.
func (h *WebhookHandler) processWebhookEventAsync(ctx context.Context, payload *ExternalEventPayload, matchID int32) {
	// Convert metadata to JSON string
	metadataJSON := ""
	if payload.Metadata != nil {
		metadataBytes, err := json.Marshal(payload.Metadata)
		if err == nil {
			metadataJSON = string(metadataBytes)
		} else {
			h.logger.Warn("Failed to marshal metadata", "error", err)
		}
	}

	// Convert float64 pointers to pgtype.Numeric
	var posX, posY pgtype.Numeric
	if payload.PositionX != nil {
		if scanErr := posX.Scan(*payload.PositionX); scanErr != nil {
			h.logger.Warn("Failed to scan PositionX", "error", scanErr, "value", *payload.PositionX)
		}
	}
	if payload.PositionY != nil {
		if scanErr := posY.Scan(*payload.PositionY); scanErr != nil {
			h.logger.Warn("Failed to scan PositionY", "error", scanErr, "value", *payload.PositionY)
		}
	}

	// Convert description to pointer
	var desc *string
	if payload.Description != "" {
		desc = &payload.Description
	}

	// Normalize event type (external: "GOAL" -> internal: "goal")
	eventType := strings.ToLower(payload.EventType)

	// Create event in database
	event, err := h.queries.CreateMatchEvent(ctx, sqlc.CreateMatchEventParams{
		MatchID:           matchID,
		TeamID:            payload.TeamID,
		PlayerID:          payload.PlayerID,
		SecondaryPlayerID: payload.SecondaryPlayerID,
		EventType:         eventType,
		Minute:            payload.Minute,
		ExtraMinute:       payload.ExtraMinute,
		PositionX:         posX,
		PositionY:         posY,
		Description:       desc,
		Metadata:          []byte(metadataJSON),
	})
	if err != nil {
		h.logger.Error("Failed to create match event from webhook", "error", err, "match_id", matchID)
		return
	}

	// Publish to real-time system (Redis Streams + Pub/Sub)
	extraMin := 0
	if payload.ExtraMinute != nil {
		extraMin = int(*payload.ExtraMinute)
	}

	var posXFloat, posYFloat *float64
	if event.PositionX.Valid {
		val, valErr := event.PositionX.Float64Value()
		if valErr == nil {
			f64 := val.Float64
			posXFloat = &f64
		}
	}
	if event.PositionY.Valid {
		val, valErr := event.PositionY.Float64Value()
		if valErr == nil {
			f64 := val.Float64
			posYFloat = &f64
		}
	}

	description := ""
	if event.Description != nil {
		description = *event.Description
	}

	// Publish to Redis Streams and Pub/Sub for real-time delivery
	publishErr := h.publisher.PublishMatchEvent(ctx, &events.MatchEvent{
		ID:                event.ID,
		MatchID:           event.MatchID,
		TeamID:            event.TeamID,
		PlayerID:          event.PlayerID,
		SecondaryPlayerID: event.SecondaryPlayerID,
		EventType:         event.EventType,
		Minute:            int(event.Minute),
		ExtraMinute:       extraMin,
		PositionX:         posXFloat,
		PositionY:         posYFloat,
		Description:       description,
		Metadata:          metadataJSON,
		Timestamp:         event.CreatedAt.Time,
	})
	if publishErr != nil {
		h.logger.Error("Failed to publish webhook event", "error", publishErr, "event_id", event.ID)
		return
	}

	// If it's a goal, invalidate match cache
	if eventType == "goal" {
		if invalidateErr := h.publisher.InvalidateMatchCache(ctx, matchID); invalidateErr != nil {
			h.logger.Warn("Failed to invalidate match cache", "error", invalidateErr, "match_id", matchID)
		}
	}

	h.logger.Info("Processed webhook event",
		"match_id", matchID,
		"event_type", eventType,
		"event_id", event.ID,
		"source", "webhook",
	)
}

// processProviderEventsAsync processes a batch of events asynchronously.
// This includes storing in DB and publishing to Redis Streams/Pub/Sub for each event.
func (h *WebhookHandler) processProviderEventsAsync(ctx context.Context, events []*events.MatchEvent, providerName string) {
	successCount := 0
	failureCount := 0
	matchIDsToInvalidate := make(map[int32]bool)

	for i, event := range events {
		// Validate match exists (should already be validated, but double-check)
		match, err := h.queries.GetMatchByID(ctx, event.MatchID)
		if err != nil {
			h.logger.Error("Match not found for batch event", "error", err, "match_id", event.MatchID, "index", i, "provider", providerName)
			failureCount++
			continue
		}

		// Process single event
		if err := h.processSingleEvent(ctx, event, match.ID, providerName); err != nil {
			h.logger.Error("Failed to process batch event", "error", err, "match_id", event.MatchID, "index", i, "provider", providerName)
			failureCount++
			continue
		}

		successCount++

		// Track matches with goals for cache invalidation
		if event.EventType == "goal" {
			matchIDsToInvalidate[event.MatchID] = true
		}
	}

	// Invalidate cache for matches with goals
	for matchID := range matchIDsToInvalidate {
		if invalidateErr := h.publisher.InvalidateMatchCache(ctx, matchID); invalidateErr != nil {
			h.logger.Warn("Failed to invalidate match cache", "error", invalidateErr, "match_id", matchID)
		}
	}

	h.logger.Info("Processed batch webhook events",
		"total", len(events),
		"success", successCount,
		"failed", failureCount,
		"provider", providerName,
	)
}

// processSingleEvent processes a single event (used by both single and batch processing).
func (h *WebhookHandler) processSingleEvent(ctx context.Context, event *events.MatchEvent, matchID int32, providerName string) error {
	// Convert metadata to JSON string if it's a map
	metadataJSON := event.Metadata
	if event.Metadata == "" {
		metadataJSON = ""
	}

	// Convert float64 pointers to pgtype.Numeric
	var posX, posY pgtype.Numeric
	if event.PositionX != nil {
		if scanErr := posX.Scan(*event.PositionX); scanErr != nil {
			h.logger.Warn("Failed to scan PositionX", "error", scanErr, "value", *event.PositionX)
		}
	}
	if event.PositionY != nil {
		if scanErr := posY.Scan(*event.PositionY); scanErr != nil {
			h.logger.Warn("Failed to scan PositionY", "error", scanErr, "value", *event.PositionY)
		}
	}

	// Convert description to pointer
	var desc *string
	if event.Description != "" {
		desc = &event.Description
	}

	// Convert second to int32 pointer
	var second *int32
	if event.Second != nil {
		s := int32(*event.Second)
		second = &s
	}

	// Convert period to string pointer
	var period *string
	if event.Period != "" {
		period = &event.Period
	}

	// Create event in database
	dbEvent, err := h.queries.CreateMatchEvent(ctx, sqlc.CreateMatchEventParams{
		MatchID:           matchID,
		TeamID:            event.TeamID,
		PlayerID:          event.PlayerID,
		SecondaryPlayerID: event.SecondaryPlayerID,
		EventType:         event.EventType,
		Minute:            int32(event.Minute),
		Second:            second,
		Period:            period,
		ExtraMinute: func() *int32 {
			if event.ExtraMinute > 0 {
				em := int32(event.ExtraMinute)
				return &em
			}
			return nil
		}(),
		PositionX:   posX,
		PositionY:   posY,
		Description: desc,
		Metadata:    []byte(metadataJSON),
	})
	if err != nil {
		return fmt.Errorf("failed to create match event: %w", err)
	}

	// Update event ID and publish to real-time system
	event.ID = dbEvent.ID
	event.Timestamp = dbEvent.CreatedAt.Time

	if publishErr := h.publisher.PublishMatchEvent(ctx, event); publishErr != nil {
		return fmt.Errorf("failed to publish event: %w", publishErr)
	}

	return nil
}

// getProviderSecret returns the secret for a specific provider.
// Falls back to default secret if provider-specific secret is not set.
func (h *WebhookHandler) getProviderSecret(providerName string) string {
	if h.webhookConfig.ProviderSecrets != nil {
		if secret, exists := h.webhookConfig.ProviderSecrets[providerName]; exists && secret != "" {
			return secret
		}
	}
	return h.webhookConfig.DefaultSecret
}

// HandleMatchStatus handles POST /webhooks/matches/:id/status.
// Receives match status updates (scheduled, live, finished, etc.) from external providers.
// @Summary Receive match status updates via webhook
// @Description Receives match status updates from external providers
// @Tags webhooks
// @Accept json
// @Produce json
// @Param id path int true "Match ID"
// @Param X-Signature header string true "HMAC SHA256 signature"
// @Param payload body map[string]interface{} true "Status update payload"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /webhooks/matches/{id}/status [post]
func (h *WebhookHandler) HandleMatchStatus(c *gin.Context) {
	// Read body for signature verification
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Warn("Failed to read request body", "error", err, "ip", c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read payload"})
		return
	}
	c.Request.Body = io.NopCloser(strings.NewReader(string(body)))

	// Verify signature using default secret
	signature := c.GetHeader("X-Signature")
	secret := h.webhookConfig.DefaultSecret
	if secret != "" {
		if !webhooks.VerifyHMACSignature(body, signature, secret) {
			h.logger.Warn("Invalid webhook signature for status update", "ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			return
		}
	}

	// Parse match ID
	matchIDStr := c.Param("id")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	// Parse payload
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Extract status
	status, ok := payload["status"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid status field"})
		return
	}

	// Validate match exists
	_, err = h.queries.GetMatchByID(c.Request.Context(), int32(matchID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	// Publish status update asynchronously
	go func() {
		statusUpdate := &events.MatchStatusUpdate{
			MatchID:   int32(matchID),
			Status:    strings.ToLower(status),
			Timestamp: time.Now(),
		}

		if publishErr := h.publisher.PublishMatchStatusUpdate(c.Request.Context(), statusUpdate); publishErr != nil {
			h.logger.Error("Failed to publish status update", "error", publishErr, "match_id", matchID)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"status":    "accepted",
		"match_id":  matchID,
		"new_status": status,
	})
}
