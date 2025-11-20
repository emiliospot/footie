package models

import (
	"time"

	"gorm.io/gorm"
)

// Player represents a football player.
type Player struct {
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DateOfBirth   *time.Time     `json:"date_of_birth,omitempty"`
	Team          *Team          `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Nationality   string         `gorm:"index" json:"nationality,omitempty"`
	FullName      string         `gorm:"not null;index" json:"full_name"`
	LastName      string         `gorm:"not null" json:"last_name"`
	Position      string         `gorm:"not null;index" json:"position"`
	PreferredFoot string         `json:"preferred_foot,omitempty"`
	Photo         string         `json:"photo,omitempty"`
	FirstName     string         `gorm:"not null" json:"first_name"`
	ID            uint           `gorm:"primarykey" json:"id"`
	Height        int            `json:"height,omitempty"`
	Weight        int            `json:"weight,omitempty"`
	ShirtNumber   int            `json:"shirt_number,omitempty"`
	TeamID        uint           `gorm:"not null;index" json:"team_id"`
}

// TableName specifies the table name for Player model.
func (Player) TableName() string {
	return "players"
}

// Age calculates the player's age.
func (p *Player) Age() int {
	if p.DateOfBirth == nil {
		return 0
	}
	return int(time.Since(*p.DateOfBirth).Hours() / 24 / 365)
}
