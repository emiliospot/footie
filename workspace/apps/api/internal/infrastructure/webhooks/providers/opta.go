package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/emiliospot/footie/api/internal/domain/events"
	infraEvents "github.com/emiliospot/footie/api/internal/infrastructure/events"
)

// OptaProvider handles Opta Sports data feed format.
// Opta uses a nested structure with event qualifiers and coordinates.
type OptaProvider struct{}

// NewOptaProvider creates a new Opta provider.
func NewOptaProvider() *OptaProvider {
	return &OptaProvider{}
}

// Name returns the provider identifier.
func (p *OptaProvider) Name() string {
	return "opta"
}

// OptaPayload represents the Opta webhook payload structure.
type OptaPayload struct {
	Event struct {
		ID        string `json:"id"`
		Type      string `json:"type"`      // e.g., "goal", "shot", "pass"
		Minute    int    `json:"minute"`
		Second    int    `json:"second,omitempty"`
		Period    string `json:"period,omitempty"` // "1H", "2H", "ET1", etc.
		MatchID   string `json:"matchId"`
		TeamID    string `json:"teamId"`
		PlayerID  string `json:"playerId,omitempty"`
		Qualifiers []struct {
			Type  string      `json:"type"`
			Value interface{} `json:"value"`
		} `json:"qualifiers,omitempty"`
		Coordinates struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		} `json:"coordinates,omitempty"`
		Description string `json:"description,omitempty"`
	} `json:"event"`
	Match struct {
		ID string `json:"id"`
	} `json:"match"`
}

// ExtractEvent extracts and transforms an Opta payload into our internal format.
func (p *OptaProvider) ExtractEvent(ctx context.Context, payload []byte) (*infraEvents.MatchEvent, error) {
	var optaPayload OptaPayload
	if err := json.Unmarshal(payload, &optaPayload); err != nil {
		return nil, fmt.Errorf("failed to parse Opta payload: %w", err)
	}

	// Convert Opta match ID to int32
	matchID, err := p.parseID(optaPayload.Match.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid match ID: %w", err)
	}

	// Convert team ID
	teamID, err := p.parseID(optaPayload.Event.TeamID)
	if err != nil {
		return nil, fmt.Errorf("invalid team ID: %w", err)
	}

	// Convert player ID (optional)
	var playerID *int32
	if optaPayload.Event.PlayerID != "" {
		pid, err := p.parseID(optaPayload.Event.PlayerID)
		if err == nil {
			playerID = &pid
		}
	}

	// Extract coordinates
	var posX, posY *float64
	if optaPayload.Event.Coordinates.X != 0 || optaPayload.Event.Coordinates.Y != 0 {
		posX = &optaPayload.Event.Coordinates.X
		posY = &optaPayload.Event.Coordinates.Y
	}

	// Extract metadata from qualifiers (xG, pass outcome, etc.)
	metadata := make(map[string]interface{})
	for _, qualifier := range optaPayload.Event.Qualifiers {
		metadata[qualifier.Type] = qualifier.Value
	}

	// Normalize and validate event type
	eventType := events.Normalize(optaPayload.Event.Type)
	if !events.IsValid(eventType) {
		return nil, fmt.Errorf("invalid event type: %s", optaPayload.Event.Type)
	}

	// Convert metadata to JSON string
	metadataJSON := ""
	if len(metadata) > 0 {
		metadataBytes, err := json.Marshal(metadata)
		if err == nil {
			metadataJSON = string(metadataBytes)
		}
	}

	// Normalize period from Opta format ("1H", "2H", "ET1", "ET2", "P")
	period := events.NormalizePeriod(optaPayload.Event.Period)
	if period == events.PeriodRegular {
		// Auto-determine if not provided
		period = events.DeterminePeriod(int32(optaPayload.Event.Minute), nil)
	}

	// Extract second (0-59)
	var second *int32
	if optaPayload.Event.Second > 0 {
		s := int32(optaPayload.Event.Second)
		if s < 60 {
			second = &s
		}
	}

	// Calculate extra minute from period
	var extraMinute *int32
	if period == events.PeriodExtraTimeFirst || period == events.PeriodExtraTimeSecond {
		extraMin := int32(0)
		extraMinute = &extraMin
	}

	return &infraEvents.MatchEvent{
		ID:                0,
		MatchID:           matchID,
		TeamID:            &teamID,
		PlayerID:          playerID,
		SecondaryPlayerID: nil,
		EventType:         eventType.String(),
		Minute:            optaPayload.Event.Minute,
		Second: func() *int {
			if second != nil {
				s := int(*second)
				return &s
			}
			return nil
		}(),
		Period: period.String(),
		ExtraMinute: func() int {
			if extraMinute != nil {
				return int(*extraMinute)
			}
			return 0
		}(),
		PositionX:   posX,
		PositionY:   posY,
		Description: optaPayload.Event.Description,
		Metadata:    metadataJSON,
	}, nil
}

// VerifySignature verifies Opta's signature format (if they use one).
func (p *OptaProvider) VerifySignature(payload []byte, signature string, secret string) bool {
	// Opta may use a different signature format
	// Implement Opta-specific verification here
	if secret == "" {
		return true
	}
	// TODO: Implement Opta-specific signature verification
	return true
}

// parseID converts Opta's string IDs to int32.
func (p *OptaProvider) parseID(idStr string) (int32, error) {
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(id), nil
}
