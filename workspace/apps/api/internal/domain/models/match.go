package models

import (
	"time"
)

// Match represents a football match.
// This is a domain model - database-agnostic, contains business logic.
type Match struct {
	ID        int32     `json:"id"`
	MatchDate time.Time `json:"match_date"`
	Competition string  `json:"competition"`
	Season      string  `json:"season"`
	Round       *string `json:"round,omitempty"`
	Stadium     *string `json:"stadium,omitempty"`
	Attendance  *int32  `json:"attendance,omitempty"`
	Status      string  `json:"status"` // scheduled, live, finished, postponed, canceled
	Referee     *string `json:"referee,omitempty"`

	// Home Team
	HomeTeamID    int32 `json:"home_team_id"`
	HomeTeamScore int32 `json:"home_team_score"`

	// Away Team
	AwayTeamID    int32 `json:"away_team_id"`
	AwayTeamScore int32 `json:"away_team_score"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"` // Soft delete timestamp
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
func (m *Match) Winner() int32 {
	if m.HomeTeamScore > m.AwayTeamScore {
		return m.HomeTeamID
	} else if m.AwayTeamScore > m.HomeTeamScore {
		return m.AwayTeamID
	}
	return 0 // Draw
}
