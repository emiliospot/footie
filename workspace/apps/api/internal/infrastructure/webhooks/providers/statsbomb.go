package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/emiliospot/footie/api/internal/domain/events"
	infraEvents "github.com/emiliospot/footie/api/internal/infrastructure/events"
)

// StatsBombProvider handles StatsBomb data feed format.
// StatsBomb uses a flat structure with location arrays.
type StatsBombProvider struct{}

// NewStatsBombProvider creates a new StatsBomb provider.
func NewStatsBombProvider() *StatsBombProvider {
	return &StatsBombProvider{}
}

// Name returns the provider identifier.
func (p *StatsBombProvider) Name() string {
	return "statsbomb"
}

// StatsBombPayload represents the StatsBomb webhook payload structure.
type StatsBombPayload struct {
	MatchID   string  `json:"match_id"`
	EventID   string  `json:"event_id"`
	Type      string  `json:"type"` // e.g., "Shot", "Pass", "Goal"
	Minute    int     `json:"minute"`
	Second    int     `json:"second,omitempty"`
	Period    int     `json:"period"` // 1 = first half, 2 = second half, etc.
	Team      string  `json:"team"`
	Player    string  `json:"player,omitempty"`
	Location  []float64 `json:"location,omitempty"` // [x, y] coordinates
	Outcome   string   `json:"outcome,omitempty"`
	BodyPart  string   `json:"body_part,omitempty"`
	Technique string   `json:"technique,omitempty"`
	XG        *float64 `json:"xG,omitempty"`
	PassEnd   []float64 `json:"pass_end_location,omitempty"`
}

// ExtractEvent extracts and transforms a single StatsBomb payload into our internal format.
func (p *StatsBombProvider) ExtractEvent(ctx context.Context, payload []byte) (*infraEvents.MatchEvent, error) {
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

// ExtractEvents extracts and transforms StatsBomb payloads (single or batch) into our internal format.
// StatsBomb batch format: array of StatsBombPayload objects.
func (p *StatsBombProvider) ExtractEvents(ctx context.Context, payload []byte) ([]*infraEvents.MatchEvent, error) {
	// Try to parse as array first (batch)
	var batchPayload []StatsBombPayload
	if err := json.Unmarshal(payload, &batchPayload); err == nil {
		// Successfully parsed as array - process batch
		events := make([]*infraEvents.MatchEvent, 0, len(batchPayload))
		for i, sbPayload := range batchPayload {
			event, err := p.extractSingleStatsBombEvent(&sbPayload)
			if err != nil {
				return nil, fmt.Errorf("failed to extract event at index %d: %w", i, err)
			}
			events = append(events, event)
		}
		return events, nil
	}

	// Try to parse as single event
	var sbPayload StatsBombPayload
	if err := json.Unmarshal(payload, &sbPayload); err != nil {
		return nil, fmt.Errorf("failed to parse StatsBomb payload as single event or batch: %w", err)
	}

	// Single event
	event, err := p.extractSingleStatsBombEvent(&sbPayload)
	if err != nil {
		return nil, err
	}
	return []*infraEvents.MatchEvent{event}, nil
}

// extractSingleStatsBombEvent extracts a single event from a StatsBombPayload.
func (p *StatsBombProvider) extractSingleStatsBombEvent(sbPayload *StatsBombPayload) (*infraEvents.MatchEvent, error) {
	if sbPayload == nil {
		return nil, fmt.Errorf("payload is nil")
	}

	// Convert match ID
	matchID, err := p.parseID(sbPayload.MatchID)
	if err != nil {
		return nil, fmt.Errorf("invalid match ID: %w", err)
	}

	// Convert team ID
	teamID, err := p.parseID(sbPayload.Team)
	if err != nil {
		return nil, fmt.Errorf("invalid team ID: %w", err)
	}

	// Convert player ID (optional)
	var playerID *int32
	if sbPayload.Player != "" {
		pid, err := p.parseID(sbPayload.Player)
		if err == nil {
			playerID = &pid
		}
	}

	// Extract coordinates from location array [x, y]
	var posX, posY *float64
	if len(sbPayload.Location) >= 2 {
		posX = &sbPayload.Location[0]
		posY = &sbPayload.Location[1]
	}

	// Build metadata from StatsBomb-specific fields
	metadata := make(map[string]interface{})
	if sbPayload.XG != nil {
		metadata["xG"] = *sbPayload.XG
	}
	if sbPayload.Outcome != "" {
		metadata["outcome"] = sbPayload.Outcome
	}
	if sbPayload.BodyPart != "" {
		metadata["body_part"] = sbPayload.BodyPart
	}
	if sbPayload.Technique != "" {
		metadata["technique"] = sbPayload.Technique
	}
	if len(sbPayload.PassEnd) >= 2 {
		metadata["pass_end_x"] = sbPayload.PassEnd[0]
		metadata["pass_end_y"] = sbPayload.PassEnd[1]
	}

	// Convert metadata to JSON string
	metadataJSON := ""
	if len(metadata) > 0 {
		metadataBytes, err := json.Marshal(metadata)
		if err == nil {
			metadataJSON = string(metadataBytes)
		}
	}

	// Normalize event type (Shot -> shot)
	eventType := events.Normalize(sbPayload.Type)
	if !events.IsValid(eventType) {
		return nil, fmt.Errorf("invalid event type: %s", sbPayload.Type)
	}

	// Convert StatsBomb period (1, 2, 3, 4, 5) to our period format
	var period events.Period
	switch sbPayload.Period {
	case 1:
		period = events.PeriodFirstHalf
	case 2:
		period = events.PeriodSecondHalf
	case 3:
		period = events.PeriodExtraTimeFirst
	case 4:
		period = events.PeriodExtraTimeSecond
	case 5:
		period = events.PeriodPenalties
	default:
		// Auto-determine from minute
		period = events.DeterminePeriod(int32(sbPayload.Minute), nil)
	}

	// Extract second (0-59)
	var second *int32
	if sbPayload.Second > 0 {
		s := int32(sbPayload.Second)
		if s < 60 {
			second = &s
		}
	}

	// Calculate extra minute (StatsBomb uses period > 2 for extra time)
	var extraMinute *int32
	if sbPayload.Period > 2 && sbPayload.Period < 5 {
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
		Minute:            sbPayload.Minute,
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
		Description: "",
		Metadata:    metadataJSON,
	}, nil
}

// VerifySignature verifies StatsBomb's signature format.
func (p *StatsBombProvider) VerifySignature(payload []byte, signature string, secret string) bool {
	// StatsBomb may use a different signature format
	if secret == "" {
		return true
	}
	// TODO: Implement StatsBomb-specific signature verification
	return true
}

// parseID converts StatsBomb's string IDs to int32.
func (p *StatsBombProvider) parseID(idStr string) (int32, error) {
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(id), nil
}
