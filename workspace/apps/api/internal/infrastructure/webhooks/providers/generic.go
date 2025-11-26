package providers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/emiliospot/footie/api/internal/domain/events"
	infraEvents "github.com/emiliospot/footie/api/internal/infrastructure/events"
	"github.com/emiliospot/footie/api/internal/infrastructure/webhooks"
)

// GenericProvider handles a simple, standardized webhook format.
// This is useful for custom integrations or providers that follow a common format.
type GenericProvider struct{}

// NewGenericProvider creates a new generic provider.
func NewGenericProvider() *GenericProvider {
	return &GenericProvider{}
}

// Name returns the provider identifier.
func (p *GenericProvider) Name() string {
	return "generic"
}

// GenericPayload represents the expected format for generic webhooks.
type GenericPayload struct {
	MatchID           int32                  `json:"matchId"`
	EventType         string                 `json:"eventType"` // GOAL, SHOT, PASS, etc.
	Minute            int32                   `json:"minute"`
	Second            *int32                  `json:"second,omitempty"` // Exact second (0-59)
	Period            *string                 `json:"period,omitempty"` // "first_half", "second_half", "extra_time_first", "extra_time_second", "penalties"
	ExtraMinute       *int32                  `json:"extraMinute,omitempty"`
	TeamID            *int32                  `json:"teamId,omitempty"`
	PlayerID          *int32                  `json:"playerId,omitempty"`
	SecondaryPlayerID *int32                  `json:"secondaryPlayerId,omitempty"`
	PositionX         *float64                `json:"positionX,omitempty"`
	PositionY         *float64                `json:"positionY,omitempty"`
	Description       string                 `json:"description,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// ExtractEvent extracts and transforms a single generic payload into our internal format.
func (p *GenericProvider) ExtractEvent(ctx context.Context, payload []byte) (*infraEvents.MatchEvent, error) {
	events, err := p.ExtractEvents(ctx, payload)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, fmt.Errorf("no events extracted from payload")
	}
	if len(events) > 1 {
		return nil, fmt.Errorf("expected single event, got %d events. Use ExtractEvents for batch processing", len(events))
	}
	return events[0], nil
}

// ExtractEvents extracts and transforms a generic payload (single or batch) into our internal format.
// Supports both single event objects and arrays of events.
func (p *GenericProvider) ExtractEvents(ctx context.Context, payload []byte) ([]*infraEvents.MatchEvent, error) {
	// Try to parse as array first (batch)
	var batchPayload []GenericPayload
	if err := json.Unmarshal(payload, &batchPayload); err == nil {
		// Successfully parsed as array - process batch
		events := make([]*infraEvents.MatchEvent, 0, len(batchPayload))
		for i, genericPayload := range batchPayload {
			event, err := p.extractSingleEvent(&genericPayload)
			if err != nil {
				return nil, fmt.Errorf("failed to extract event at index %d: %w", i, err)
			}
			events = append(events, event)
		}
		return events, nil
	}

	// Try to parse as single event
	var genericPayload GenericPayload
	if err := json.Unmarshal(payload, &genericPayload); err != nil {
		return nil, fmt.Errorf("failed to parse payload as single event or batch: %w", err)
	}

	// Single event
	event, err := p.extractSingleEvent(&genericPayload)
	if err != nil {
		return nil, err
	}
	return []*infraEvents.MatchEvent{event}, nil
}

// extractSingleEvent extracts a single event from a GenericPayload.
func (p *GenericProvider) extractSingleEvent(genericPayload *GenericPayload) (*infraEvents.MatchEvent, error) {
	// Normalize and validate event type (GOAL -> goal)
	normalizedType := events.Normalize(genericPayload.EventType)
	if !events.IsValid(normalizedType) {
		return nil, fmt.Errorf("invalid event type: %s", genericPayload.EventType)
	}

	// Normalize period (if provided, otherwise determine from minute)
	var period string
	if genericPayload.Period != nil {
		normalizedPeriod := events.NormalizePeriod(*genericPayload.Period)
		period = normalizedPeriod.String()
	} else {
		// Auto-determine period from minute and extra_minute
		period = events.DeterminePeriod(genericPayload.Minute, genericPayload.ExtraMinute).String()
	}

	// Validate second (0-59)
	var second *int32
	if genericPayload.Second != nil {
		s := *genericPayload.Second
		if s < 0 || s >= 60 {
			return nil, fmt.Errorf("invalid second: %d (must be 0-59)", s)
		}
		second = &s
	}

	// Convert metadata to JSON string
	metadataJSON := ""
	if genericPayload.Metadata != nil {
		metadataBytes, err := json.Marshal(genericPayload.Metadata)
		if err == nil {
			metadataJSON = string(metadataBytes)
		}
	}

	// Convert to infraEvents.MatchEvent (temporary ID, will be set after DB insert)
	return &infraEvents.MatchEvent{
		ID:                0, // Will be set after DB insert
		MatchID:           genericPayload.MatchID,
		TeamID:            genericPayload.TeamID,
		PlayerID:          genericPayload.PlayerID,
		SecondaryPlayerID: genericPayload.SecondaryPlayerID,
		EventType:         normalizedType.String(),
		Minute:            int(genericPayload.Minute),
		Second: func() *int {
			if second != nil {
				s := int(*second)
				return &s
			}
			return nil
		}(),
		Period: period,
		ExtraMinute: func() int {
			if genericPayload.ExtraMinute != nil {
				return int(*genericPayload.ExtraMinute)
			}
			return 0
		}(),
		PositionX:   genericPayload.PositionX,
		PositionY:   genericPayload.PositionY,
		Description: genericPayload.Description,
		Metadata:    metadataJSON,
	}, nil
}

// VerifySignature verifies the HMAC SHA256 signature (standard implementation).
func (p *GenericProvider) VerifySignature(payload []byte, signature string, secret string) bool {
	return webhooks.VerifyHMACSignature(payload, signature, secret)
}
