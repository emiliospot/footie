package webhooks

import (
	"context"

	"github.com/emiliospot/footie/api/internal/infrastructure/events"
)

// Provider defines the interface for external data feed providers.
// Each provider (Opta, StatsBomb, API-Football, etc.) implements this interface
// to transform their specific payload format into our internal event format.
type Provider interface {
	// Name returns the provider identifier (e.g., "opta", "statsbomb", "api-football").
	Name() string

	// ExtractEvent extracts and transforms a raw JSON payload into our internal MatchEvent format.
	// The payload is the raw JSON bytes from the webhook request.
	// This method handles single events.
	ExtractEvent(ctx context.Context, payload []byte) (*events.MatchEvent, error)

	// ExtractEvents extracts and transforms a batch of events from a raw JSON payload.
	// The payload can be either a single event object or an array of events.
	// Returns an array of MatchEvent, which will have length 1 for single events.
	// Providers should implement this to support batch processing.
	ExtractEvents(ctx context.Context, payload []byte) ([]*events.MatchEvent, error)

	// VerifySignature verifies the webhook signature for this provider.
	// Returns true if the signature is valid, false otherwise.
	// Some providers may not use signatures (returns true in that case).
	VerifySignature(payload []byte, signature string, secret string) bool
}

// NormalizedEvent represents the internal event format that all providers must produce.
// This is what gets stored in the database and published to Redis.
type NormalizedEvent struct {
	MatchID           int32
	EventType         string // normalized: "goal", "shot", "pass", etc.
	Minute            int32
	ExtraMinute       *int32
	TeamID            *int32
	PlayerID          *int32
	SecondaryPlayerID *int32
	PositionX         *float64
	PositionY         *float64
	Description       string
	Metadata          map[string]interface{} // provider-specific metadata (xG, pass completion, etc.)
}
