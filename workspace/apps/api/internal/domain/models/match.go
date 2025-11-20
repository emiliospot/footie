package models

import (
	"time"

	"gorm.io/gorm"
)

// Match represents a football match.
//
//nolint:govet // Field alignment optimization would reduce readability
type Match struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	MatchDate   time.Time `gorm:"not null;index" json:"match_date"`
	Competition string    `gorm:"not null;index" json:"competition"`
	Season      string    `gorm:"not null;index" json:"season"`
	Round       string    `json:"round,omitempty"`
	Stadium     string    `json:"stadium,omitempty"`
	Attendance  int       `json:"attendance,omitempty"`
	Status      string    `gorm:"not null;default:'scheduled'" json:"status"` // scheduled, live, finished, postponed, canceled
	Referee     string    `json:"referee,omitempty"`

	// Home Team
	HomeTeamID    uint  `gorm:"not null;index" json:"home_team_id"`
	HomeTeam      *Team `gorm:"foreignKey:HomeTeamID" json:"home_team,omitempty"`
	HomeTeamScore int   `gorm:"default:0" json:"home_team_score"`

	// Away Team
	AwayTeamID    uint  `gorm:"not null;index" json:"away_team_id"`
	AwayTeam      *Team `gorm:"foreignKey:AwayTeamID" json:"away_team,omitempty"`
	AwayTeamScore int   `gorm:"default:0" json:"away_team_score"`

	// Relations
	Events []MatchEvent `gorm:"foreignKey:MatchID" json:"events,omitempty"`
}

// TableName specifies the table name for Match model.
func (Match) TableName() string {
	return "matches"
}

// IsFinished returns true if the match is finished.
func (m *Match) IsFinished() bool {
	return m.Status == "finished"
}

// IsLive returns true if the match is currently live.
func (m *Match) IsLive() bool {
	return m.Status == "live"
}

// Winner returns the ID of the winning team, or 0 for a draw.
func (m *Match) Winner() uint {
	if m.HomeTeamScore > m.AwayTeamScore {
		return m.HomeTeamID
	} else if m.AwayTeamScore > m.HomeTeamScore {
		return m.AwayTeamID
	}
	return 0 // Draw
}
