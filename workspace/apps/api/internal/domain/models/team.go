package models

import (
	"time"

	"gorm.io/gorm"
)

// Team represents a football team.
type Team struct {
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Country         string         `gorm:"not null;index" json:"country"`
	Name            string         `gorm:"not null;index" json:"name"`
	ShortName       string         `gorm:"not null" json:"short_name"`
	Code            string         `gorm:"uniqueIndex;not null" json:"code"`
	City            string         `json:"city,omitempty"`
	Stadium         string         `json:"stadium,omitempty"`
	Logo            string         `json:"logo,omitempty"`
	Colors          string         `json:"colors,omitempty"`
	Website         string         `json:"website,omitempty"`
	Players         []Player       `gorm:"foreignKey:TeamID" json:"players,omitempty"`
	Founded         int            `json:"founded,omitempty"`
	ID              uint           `gorm:"primarykey" json:"id"`
	StadiumCapacity int            `json:"stadium_capacity,omitempty"`
}

// TableName specifies the table name for Team model.
func (Team) TableName() string {
	return "teams"
}
