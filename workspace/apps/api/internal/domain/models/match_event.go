package models

import (
	"time"

	"gorm.io/gorm"
)

// MatchEvent represents an event that occurred during a match.
type MatchEvent struct {
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	PlayerID          *uint          `gorm:"index" json:"player_id,omitempty"`
	SecondaryPlayer   *Player        `gorm:"foreignKey:SecondaryPlayerID" json:"secondary_player,omitempty"`
	Match             *Match         `gorm:"foreignKey:MatchID" json:"match,omitempty"`
	SecondaryPlayerID *uint          `gorm:"index" json:"secondary_player_id,omitempty"`
	Team              *Team          `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	TeamID            *uint          `gorm:"index" json:"team_id,omitempty"`
	Player            *Player        `gorm:"foreignKey:PlayerID" json:"player,omitempty"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	EventType         string         `gorm:"not null;index" json:"event_type"`
	Description       string         `json:"description,omitempty"`
	Metadata          string         `gorm:"type:jsonb" json:"metadata,omitempty"`
	ID                uint           `gorm:"primarykey" json:"id"`
	ExtraMinute       int            `json:"extra_minute,omitempty"`
	Minute            int            `gorm:"not null" json:"minute"`
	MatchID           uint           `gorm:"not null;index" json:"match_id"`
	PositionX         float64        `json:"position_x,omitempty"`
	PositionY         float64        `json:"position_y,omitempty"`
}

// TableName specifies the table name for MatchEvent model.
func (MatchEvent) TableName() string {
	return "match_events"
}

// IsGoal returns true if the event is a goal.
func (me *MatchEvent) IsGoal() bool {
	return me.EventType == "goal"
}

// IsCard returns true if the event is a card (yellow or red).
func (me *MatchEvent) IsCard() bool {
	return me.EventType == "yellow_card" || me.EventType == "red_card"
}
