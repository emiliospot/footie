package models

import (
	"time"
)

// Player represents a football player.
// This is a domain model - database-agnostic, contains business logic.
type Player struct {
	ID            int32      `json:"id"`
	TeamID        int32      `json:"team_id"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	FullName      string     `json:"full_name"`
	DateOfBirth   *time.Time `json:"date_of_birth,omitempty"`
	Nationality   *string    `json:"nationality,omitempty"`
	Position      string     `json:"position"`
	ShirtNumber   *int32     `json:"shirt_number,omitempty"`
	Height        *int32     `json:"height,omitempty"` // in cm
	Weight        *int32     `json:"weight,omitempty"` // in kg
	PreferredFoot *string    `json:"preferred_foot,omitempty"`
	Photo         *string    `json:"photo,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"-"` // Soft delete timestamp
}

// Age calculates the player's age.
func (p *Player) Age() int {
	if p.DateOfBirth == nil {
		return 0
	}
	return int(time.Since(*p.DateOfBirth).Hours() / 24 / 365)
}
