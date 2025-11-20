package models

import (
	"time"

	"gorm.io/gorm"
)

// PlayerStatistics represents aggregated statistics for a player.
type PlayerStatistics struct {
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	Player          *Player        `gorm:"foreignKey:PlayerID" json:"player,omitempty"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Season          string         `gorm:"not null;index" json:"season"`
	Competition     string         `gorm:"not null;index" json:"competition"`
	PassesCompleted int            `json:"passes_completed"`
	Crosses         int            `json:"crosses"`
	MatchesPlayed   int            `json:"matches_played"`
	MinutesPlayed   int            `json:"minutes_played"`
	MatchesStarted  int            `json:"matches_started"`
	SubOn           int            `json:"sub_on"`
	SubOff          int            `json:"sub_off"`
	Goals           int            `json:"goals"`
	Assists         int            `json:"assists"`
	ShotsTotal      int            `json:"shots_total"`
	ShotsOnTarget   int            `json:"shots_on_target"`
	ShotAccuracy    float64        `json:"shot_accuracy"`
	GoalConversion  float64        `json:"goal_conversion"`
	PassesTotal     int            `json:"passes_total"`
	ID              uint           `gorm:"primarykey" json:"id"`
	PassAccuracy    float64        `json:"pass_accuracy"`
	KeyPasses       int            `json:"key_passes"`
	PlayerID        uint           `gorm:"not null;index" json:"player_id"`
	Tackles         int            `json:"tackles"`
	TacklesWon      int            `json:"tackles_won"`
	Interceptions   int            `json:"interceptions"`
	Clearances      int            `json:"clearances"`
	BlockedShots    int            `json:"blocked_shots"`
	Duels           int            `json:"duels"`
	DuelsWon        int            `json:"duels_won"`
	AerialDuels     int            `json:"aerial_duels"`
	AerialDuelsWon  int            `json:"aerial_duels_won"`
	YellowCards     int            `json:"yellow_cards"`
	RedCards        int            `json:"red_cards"`
	Fouls           int            `json:"fouls"`
	FoulsDrawn      int            `json:"fouls_drawn"`
	CleanSheets     int            `json:"clean_sheets,omitempty"`
	GoalsConceded   int            `json:"goals_conceded,omitempty"`
	SavesTotal      int            `json:"saves_total,omitempty"`
	SavePercentage  float64        `json:"save_percentage,omitempty"`
	PenaltiesSaved  int            `json:"penalties_saved,omitempty"`
}

// TableName specifies the table name for PlayerStatistics model.
func (PlayerStatistics) TableName() string {
	return "player_statistics"
}

// TeamStatistics represents aggregated statistics for a team.
type TeamStatistics struct {
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Team           *Team          `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Competition    string         `gorm:"not null;index" json:"competition"`
	CurrentForm    string         `json:"current_form"`
	Season         string         `gorm:"not null;index" json:"season"`
	GoalsPerMatch  float64        `json:"goals_per_match"`
	ShotsPerMatch  float64        `json:"shots_per_match"`
	Wins           int            `json:"wins"`
	Draws          int            `json:"draws"`
	Losses         int            `json:"losses"`
	Points         int            `json:"points"`
	GoalsScored    int            `json:"goals_scored"`
	GoalsConceded  int            `json:"goals_conceded"`
	GoalDifference int            `json:"goal_difference"`
	ID             uint           `gorm:"primarykey" json:"id"`
	CleanSheets    int            `json:"clean_sheets"`
	Possession     float64        `json:"possession"`
	PassAccuracy   float64        `json:"pass_accuracy"`
	MatchesPlayed  int            `json:"matches_played"`
	ShotsOnTarget  float64        `json:"shots_on_target_percentage"`
	HomeWins       int            `json:"home_wins"`
	HomeDraws      int            `json:"home_draws"`
	HomeLosses     int            `json:"home_losses"`
	AwayWins       int            `json:"away_wins"`
	AwayDraws      int            `json:"away_draws"`
	AwayLosses     int            `json:"away_losses"`
	YellowCards    int            `json:"yellow_cards"`
	RedCards       int            `json:"red_cards"`
	TeamID         uint           `gorm:"not null;index" json:"team_id"`
	Position       int            `json:"position,omitempty"`
}

// TableName specifies the table name for TeamStatistics model.
func (TeamStatistics) TableName() string {
	return "team_statistics"
}

// WinPercentage calculates the win percentage.
func (ts *TeamStatistics) WinPercentage() float64 {
	if ts.MatchesPlayed == 0 {
		return 0
	}
	return float64(ts.Wins) / float64(ts.MatchesPlayed) * 100
}
