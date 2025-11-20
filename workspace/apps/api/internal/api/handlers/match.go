package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/emiliospot/footie/api/internal/infrastructure/events"
	"github.com/emiliospot/footie/api/internal/repository/sqlc"
)

// MatchHandler handles match-related endpoints.
type MatchHandler struct {
	*BaseHandler
}

// NewMatchHandler creates a new match handler.
func NewMatchHandler(base *BaseHandler) *MatchHandler {
	return &MatchHandler{BaseHandler: base}
}

// ListMatchesRequest represents the query parameters for listing matches.
type ListMatchesRequest struct {
	Limit  int32 `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset int32 `form:"offset" binding:"omitempty,min=0"`
}

// CreateMatchEventRequest represents the request to create a match event.
type CreateMatchEventRequest struct {
	TeamID            *int32   `json:"team_id"`
	PlayerID          *int32   `json:"player_id"`
	SecondaryPlayerID *int32   `json:"secondary_player_id"`
	EventType         string   `json:"event_type" binding:"required"`
	Minute            int32    `json:"minute" binding:"required,min=0,max=120"`
	ExtraMinute       *int32   `json:"extra_minute"`
	PositionX         *float64 `json:"position_x"`
	PositionY         *float64 `json:"position_y"`
	Description       string   `json:"description"`
	Metadata          string   `json:"metadata"`
}

// ListMatches handles GET /api/v1/matches
// @Summary List matches
// @Description Get a list of matches
// @Tags matches
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} sqlc.Match
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/matches [get]
func (h *MatchHandler) ListMatches(c *gin.Context) {
	var req ListMatchesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if req.Limit == 0 {
		req.Limit = 20
	}

	matches, err := h.queries.ListMatches(c.Request.Context(), sqlc.ListMatchesParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		h.logger.Error("Failed to list matches", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve matches"})
		return
	}

	c.JSON(http.StatusOK, matches)
}

// GetMatch handles GET /api/v1/matches/:id
// @Summary Get match by ID
// @Description Get detailed information about a specific match
// @Tags matches
// @Accept json
// @Produce json
// @Param id path int true "Match ID"
// @Success 200 {object} sqlc.Match
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/matches/{id} [get]
func (h *MatchHandler) GetMatch(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	match, err := h.queries.GetMatchByID(c.Request.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
			return
		}
		h.logger.Error("Failed to get match", "error", err, "match_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve match"})
		return
	}

	c.JSON(http.StatusOK, match)
}

// GetMatchEvents handles GET /api/v1/matches/:id/events
// @Summary Get match events
// @Description Get all events for a specific match
// @Tags matches
// @Accept json
// @Produce json
// @Param id path int true "Match ID"
// @Param limit query int false "Limit" default(100)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} sqlc.MatchEvent
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/matches/{id}/events [get]
func (h *MatchHandler) GetMatchEvents(c *gin.Context) {
	idStr := c.Param("id")
	matchID, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	var req ListMatchesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Limit == 0 {
		req.Limit = 100
	}

	matchEvents, err := h.queries.GetMatchEvents(c.Request.Context(), int32(matchID))
	if err != nil {
		h.logger.Error("Failed to get match events", "error", err, "match_id", matchID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve match events"})
		return
	}

	c.JSON(http.StatusOK, matchEvents)
}

// CreateMatchEvent handles POST /api/v1/matches/:id/events
// @Summary Create match event
// @Description Create a new event for a match (goal, shot, pass, etc.) and broadcast it in real-time
// @Tags matches
// @Accept json
// @Produce json
// @Param id path int true "Match ID"
// @Param event body CreateMatchEventRequest true "Match event data"
// @Success 201 {object} sqlc.MatchEvent
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/matches/{id}/events [post]
func (h *MatchHandler) CreateMatchEvent(c *gin.Context) {
	idStr := c.Param("id")
	matchID, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	var req CreateMatchEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate metadata is valid JSON if provided
	if req.Metadata != "" {
		var metadataCheck map[string]interface{}
		if err := json.Unmarshal([]byte(req.Metadata), &metadataCheck); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metadata JSON"})
			return
		}
	}

	// Convert float64 pointers to pgtype.Numeric
	var posX, posY pgtype.Numeric
	if req.PositionX != nil {
		_ = posX.Scan(*req.PositionX)
	}
	if req.PositionY != nil {
		_ = posY.Scan(*req.PositionY)
	}

	// Convert description to pointer
	var desc *string
	if req.Description != "" {
		desc = &req.Description
	}

	// Create event in database
	event, err := h.queries.CreateMatchEvent(c.Request.Context(), sqlc.CreateMatchEventParams{
		MatchID:           int32(matchID),
		TeamID:            req.TeamID,
		PlayerID:          req.PlayerID,
		SecondaryPlayerID: req.SecondaryPlayerID,
		EventType:         req.EventType,
		Minute:            req.Minute,
		ExtraMinute:       req.ExtraMinute,
		PositionX:         posX,
		PositionY:         posY,
		Description:       desc,
		Metadata:          []byte(req.Metadata),
	})
	if err != nil {
		h.logger.Error("Failed to create match event", "error", err, "match_id", matchID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create match event"})
		return
	}

	// Publish to real-time system (Redis Streams + Pub/Sub)
	go func() {
		extraMin := 0
		if event.ExtraMinute != nil {
			extraMin = int(*event.ExtraMinute)
		}

		var posXFloat, posYFloat *float64
		if event.PositionX.Valid {
			val, _ := event.PositionX.Float64Value()
			f64 := val.Float64
			posXFloat = &f64
		}
		if event.PositionY.Valid {
			val, _ := event.PositionY.Float64Value()
			f64 := val.Float64
			posYFloat = &f64
		}

		description := ""
		if event.Description != nil {
			description = *event.Description
		}

		publishErr := h.publisher.PublishMatchEvent(c.Request.Context(), &events.MatchEvent{
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
			Metadata:          string(event.Metadata),
			Timestamp:         event.CreatedAt.Time,
		})
		if publishErr != nil {
			h.logger.Error("Failed to publish match event", "error", publishErr, "event_id", event.ID)
		}
	}()

	// If it's a goal, update the score
	if req.EventType == "goal" {
		// TODO: Fetch current match scores and publish score update
		// This would require additional queries to get home/away scores
		h.logger.Info("Goal scored", "match_id", matchID, "player_id", req.PlayerID)
	}

	h.logger.Info("Match event created",
		"match_id", matchID,
		"event_type", req.EventType,
		"event_id", event.ID,
	)

	c.JSON(http.StatusCreated, event)
}
