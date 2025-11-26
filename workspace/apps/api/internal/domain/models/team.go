package models

import (
	"time"
)

// Team represents a football team.
// This is a domain model - database-agnostic, contains business logic.
type Team struct {
	ID              int32     `json:"id"`
	Name            string    `json:"name"`
	ShortName       string    `json:"short_name"`
	Code            string    `json:"code"`
	Country         string    `json:"country"`
	City            *string   `json:"city,omitempty"`
	Stadium         *string   `json:"stadium,omitempty"`
	StadiumCapacity *int32    `json:"stadium_capacity,omitempty"`
	Founded         *int32    `json:"founded,omitempty"`
	Logo            *string   `json:"logo,omitempty"`
	Colors          *string   `json:"colors,omitempty"`
	Website         *string   `json:"website,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"-"` // Soft delete timestamp
}
