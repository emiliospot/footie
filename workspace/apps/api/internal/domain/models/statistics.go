package models

import (
	"time"
)

// PlayerStatistics represents aggregated statistics for a player.
// This is a domain model - database-agnostic, contains business logic.
type PlayerStatistics struct {
	ID              int32   `json:"id"`
	PlayerID        int32   `json:"player_id"`
	Season          string  `json:"season"`
	Competition     string  `json:"competition"`
	MatchesPlayed   int32   `json:"matches_played"`
	MatchesStarted  int32   `json:"matches_started"`
	MinutesPlayed   int32   `json:"minutes_played"`
	SubOn           int32   `json:"sub_on"`
	SubOff          int32   `json:"sub_off"`
	Goals           int32   `json:"goals"`
	Assists         int32   `json:"assists"`
	ShotsTotal      int32   `json:"shots_total"`
	ShotsOnTarget   int32   `json:"shots_on_target"`
	ShotAccuracy    *float64 `json:"shot_accuracy,omitempty"`
	GoalConversion  *float64 `json:"goal_conversion,omitempty"`
	PassesTotal     int32   `json:"passes_total"`
	PassesCompleted int32   `json:"passes_completed"`
	PassAccuracy    *float64 `json:"pass_accuracy,omitempty"`
	KeyPasses       int32   `json:"key_passes"`
	Crosses         int32   `json:"crosses"`
	Tackles         int32   `json:"tackles"`
	TacklesWon      int32   `json:"tackles_won"`
	Interceptions   int32   `json:"interceptions"`
	Clearances      int32   `json:"clearances"`
	BlockedShots    int32   `json:"blocked_shots"`
	Duels           int32   `json:"duels"`
	DuelsWon        int32   `json:"duels_won"`
	AerialDuels     int32   `json:"aerial_duels"`
	AerialDuelsWon  int32   `json:"aerial_duels_won"`
	YellowCards     int32   `json:"yellow_cards"`
	RedCards        int32   `json:"red_cards"`
	Fouls           int32   `json:"fouls"`
	FoulsDrawn      int32   `json:"fouls_drawn"`
	CleanSheets     *int32  `json:"clean_sheets,omitempty"`
	GoalsConceded   *int32  `json:"goals_conceded,omitempty"`
	SavesTotal      *int32  `json:"saves_total,omitempty"`
	SavePercentage  *float64 `json:"save_percentage,omitempty"`
	PenaltiesSaved  *int32  `json:"penalties_saved,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"-"` // Soft delete timestamp
}

// TeamStatistics represents aggregated statistics for a team.
// This is a domain model - database-agnostic, contains business logic.
type TeamStatistics struct {
	ID                      int32    `json:"id"`
	TeamID                  int32    `json:"team_id"`
	Season                  string   `json:"season"`
	Competition             string   `json:"competition"`
	MatchesPlayed           int32    `json:"matches_played"`
	Wins                    int32    `json:"wins"`
	Draws                   int32    `json:"draws"`
	Losses                  int32    `json:"losses"`
	Points                  int32    `json:"points"`
	Position                *int32   `json:"position,omitempty"`
	GoalsScored             int32    `json:"goals_scored"`
	GoalsConceded           int32    `json:"goals_conceded"`
	GoalDifference          int32    `json:"goal_difference"`
	CleanSheets             int32    `json:"clean_sheets"`
	GoalsPerMatch           *float64 `json:"goals_per_match,omitempty"`
	HomeWins                int32    `json:"home_wins"`
	HomeDraws               int32    `json:"home_draws"`
	HomeLosses              int32    `json:"home_losses"`
	AwayWins                int32    `json:"away_wins"`
	AwayDraws               int32    `json:"away_draws"`
	AwayLosses              int32    `json:"away_losses"`
	Possession              *float64 `json:"possession,omitempty"`
	PassAccuracy            *float64 `json:"pass_accuracy,omitempty"`
	ShotsPerMatch           *float64 `json:"shots_per_match,omitempty"`
	ShotsOnTargetPercentage *float64 `json:"shots_on_target_percentage,omitempty"`
	YellowCards             int32    `json:"yellow_cards"`
	RedCards                int32    `json:"red_cards"`
	CurrentForm             *string  `json:"current_form,omitempty"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	DeletedAt               *time.Time `json:"-"` // Soft delete timestamp
}

// WinPercentage calculates the win percentage.
func (ts *TeamStatistics) WinPercentage() float64 {
	if ts.MatchesPlayed == 0 {
		return 0
	}
	return float64(ts.Wins) / float64(ts.MatchesPlayed) * 100
}
