package models

import (
	"encoding/json"
	"time"
)

// MatchEvent represents an event that occurred during a match.
// This is a domain model - database-agnostic, contains business logic.
type MatchEvent struct {
	ID                int32           `json:"id"`
	MatchID           int32           `json:"match_id"`
	TeamID            *int32          `json:"team_id,omitempty"`
	PlayerID          *int32          `json:"player_id,omitempty"`
	SecondaryPlayerID *int32          `json:"secondary_player_id,omitempty"`
	EventType         string          `json:"event_type"`
	Minute            int32           `json:"minute"`
	ExtraMinute       *int32          `json:"extra_minute,omitempty"`
	PositionX         *float64        `json:"position_x,omitempty"`
	PositionY         *float64        `json:"position_y,omitempty"`
	Description       *string         `json:"description,omitempty"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	DeletedAt         *time.Time      `json:"-"` // Soft delete timestamp
}

// IsGoal returns true if the event is a goal.
func (me *MatchEvent) IsGoal() bool {
	return me.EventType == "goal"
}

// IsCard returns true if the event is a card (yellow or red).
func (me *MatchEvent) IsCard() bool {
	return me.EventType == "yellow_card" || me.EventType == "red_card"
}
