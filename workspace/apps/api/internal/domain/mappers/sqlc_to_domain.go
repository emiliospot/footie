package mappers

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/emiliospot/footie/api/internal/domain/models"
	"github.com/emiliospot/footie/api/internal/repository/sqlc"
)

// ToDomainUser converts a sqlc.User to a domain models.User.
func ToDomainUser(u *sqlc.User) models.User {
	return models.User{
		ID:            u.ID,
		Email:         u.Email,
		PasswordHash:  u.PasswordHash,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Role:          u.Role,
		Avatar:        u.Avatar,
		Organization:  u.Organization,
		IsActive:      u.IsActive,
		EmailVerified: u.EmailVerified,
		CreatedAt:     pgtypeToTime(u.CreatedAt),
		UpdatedAt:     pgtypeToTime(u.UpdatedAt),
		DeletedAt:     pgtypeToTimePtr(u.DeletedAt),
	}
}

// ToDomainTeam converts a sqlc.Team to a domain models.Team.
func ToDomainTeam(t *sqlc.Team) models.Team {
	return models.Team{
		ID:              t.ID,
		Name:            t.Name,
		ShortName:       t.ShortName,
		Code:            t.Code,
		Country:         t.Country,
		City:            t.City,
		Stadium:         t.Stadium,
		StadiumCapacity: t.StadiumCapacity,
		Founded:         t.Founded,
		Logo:            t.Logo,
		Colors:          t.Colors,
		Website:         t.Website,
		CreatedAt:       pgtypeToTime(t.CreatedAt),
		UpdatedAt:       pgtypeToTime(t.UpdatedAt),
		DeletedAt:       pgtypeToTimePtr(t.DeletedAt),
	}
}

// ToDomainPlayer converts a sqlc.Player to a domain models.Player.
func ToDomainPlayer(p *sqlc.Player) models.Player {
	var dateOfBirth *time.Time
	if p.DateOfBirth.Valid {
		dateOfBirth = &p.DateOfBirth.Time
	}

	return models.Player{
		ID:            p.ID,
		TeamID:        p.TeamID,
		FirstName:     p.FirstName,
		LastName:      p.LastName,
		FullName:      p.FullName,
		DateOfBirth:   dateOfBirth,
		Nationality:   p.Nationality,
		Position:      p.Position,
		ShirtNumber:   p.ShirtNumber,
		Height:        p.Height,
		Weight:        p.Weight,
		PreferredFoot: p.PreferredFoot,
		Photo:         p.Photo,
		CreatedAt:     pgtypeToTime(p.CreatedAt),
		UpdatedAt:     pgtypeToTime(p.UpdatedAt),
		DeletedAt:     pgtypeToTimePtr(p.DeletedAt),
	}
}

// ToDomainMatch converts a sqlc.Match to a domain models.Match.
func ToDomainMatch(m *sqlc.Match) models.Match {
	return models.Match{
		ID:            m.ID,
		MatchDate:     pgtypeToTime(m.MatchDate),
		Competition:   m.Competition,
		Season:        m.Season,
		Round:         m.Round,
		Stadium:       m.Stadium,
		Attendance:    m.Attendance,
		Status:        m.Status,
		Referee:       m.Referee,
		HomeTeamID:    m.HomeTeamID,
		HomeTeamScore: m.HomeTeamScore,
		AwayTeamID:    m.AwayTeamID,
		AwayTeamScore: m.AwayTeamScore,
		CreatedAt:     pgtypeToTime(m.CreatedAt),
		UpdatedAt:     pgtypeToTime(m.UpdatedAt),
		DeletedAt:     pgtypeToTimePtr(m.DeletedAt),
	}
}

// ToDomainMatchEvent converts a sqlc.MatchEvent to a domain models.MatchEvent.
func ToDomainMatchEvent(e *sqlc.MatchEvent) models.MatchEvent {
	var posX, posY *float64
	if e.PositionX.Valid {
		if val, err := e.PositionX.Float64Value(); err == nil {
			posX = &val.Float64
		}
	}
	if e.PositionY.Valid {
		if val, err := e.PositionY.Float64Value(); err == nil {
			posY = &val.Float64
		}
	}

	return models.MatchEvent{
		ID:                e.ID,
		MatchID:           e.MatchID,
		TeamID:            e.TeamID,
		PlayerID:          e.PlayerID,
		SecondaryPlayerID: e.SecondaryPlayerID,
		EventType:         e.EventType,
		Minute:            e.Minute,
		ExtraMinute:       e.ExtraMinute,
		PositionX:         posX,
		PositionY:         posY,
		Description:       e.Description,
		Metadata:          e.Metadata,
		CreatedAt:         pgtypeToTime(e.CreatedAt),
		UpdatedAt:         pgtypeToTime(e.UpdatedAt),
		DeletedAt:         pgtypeToTimePtr(e.DeletedAt),
	}
}

// ToDomainPlayerStatistics converts a sqlc.PlayerStatistic to a domain models.PlayerStatistics.
func ToDomainPlayerStatistics(s *sqlc.PlayerStatistic) models.PlayerStatistics {
	return models.PlayerStatistics{
		ID:              s.ID,
		PlayerID:        s.PlayerID,
		Season:          s.Season,
		Competition:     s.Competition,
		MatchesPlayed:   s.MatchesPlayed,
		MatchesStarted:  s.MatchesStarted,
		MinutesPlayed:   s.MinutesPlayed,
		SubOn:           s.SubOn,
		SubOff:          s.SubOff,
		Goals:           s.Goals,
		Assists:         s.Assists,
		ShotsTotal:      s.ShotsTotal,
		ShotsOnTarget:   s.ShotsOnTarget,
		ShotAccuracy:    numericToFloat64Ptr(s.ShotAccuracy),
		GoalConversion:  numericToFloat64Ptr(s.GoalConversion),
		PassesTotal:     s.PassesTotal,
		PassesCompleted: s.PassesCompleted,
		PassAccuracy:    numericToFloat64Ptr(s.PassAccuracy),
		KeyPasses:       s.KeyPasses,
		Crosses:         s.Crosses,
		Tackles:         s.Tackles,
		TacklesWon:      s.TacklesWon,
		Interceptions:   s.Interceptions,
		Clearances:      s.Clearances,
		BlockedShots:    s.BlockedShots,
		Duels:           s.Duels,
		DuelsWon:        s.DuelsWon,
		AerialDuels:     s.AerialDuels,
		AerialDuelsWon:  s.AerialDuelsWon,
		YellowCards:     s.YellowCards,
		RedCards:        s.RedCards,
		Fouls:           s.Fouls,
		FoulsDrawn:      s.FoulsDrawn,
		CleanSheets:     s.CleanSheets,
		GoalsConceded:   s.GoalsConceded,
		SavesTotal:      s.SavesTotal,
		SavePercentage:  numericToFloat64Ptr(s.SavePercentage),
		PenaltiesSaved:  s.PenaltiesSaved,
		CreatedAt:       pgtypeToTime(s.CreatedAt),
		UpdatedAt:       pgtypeToTime(s.UpdatedAt),
		DeletedAt:       pgtypeToTimePtr(s.DeletedAt),
	}
}

// ToDomainTeamStatistics converts a sqlc.TeamStatistic to a domain models.TeamStatistics.
func ToDomainTeamStatistics(s *sqlc.TeamStatistic) models.TeamStatistics {
	return models.TeamStatistics{
		ID:                      s.ID,
		TeamID:                  s.TeamID,
		Season:                  s.Season,
		Competition:             s.Competition,
		MatchesPlayed:           s.MatchesPlayed,
		Wins:                    s.Wins,
		Draws:                   s.Draws,
		Losses:                  s.Losses,
		Points:                  s.Points,
		Position:                s.Position,
		GoalsScored:             s.GoalsScored,
		GoalsConceded:           s.GoalsConceded,
		GoalDifference:          s.GoalDifference,
		CleanSheets:             s.CleanSheets,
		GoalsPerMatch:           numericToFloat64Ptr(s.GoalsPerMatch),
		HomeWins:                s.HomeWins,
		HomeDraws:               s.HomeDraws,
		HomeLosses:              s.HomeLosses,
		AwayWins:                s.AwayWins,
		AwayDraws:               s.AwayDraws,
		AwayLosses:              s.AwayLosses,
		Possession:              numericToFloat64Ptr(s.Possession),
		PassAccuracy:            numericToFloat64Ptr(s.PassAccuracy),
		ShotsPerMatch:           numericToFloat64Ptr(s.ShotsPerMatch),
		ShotsOnTargetPercentage: numericToFloat64Ptr(s.ShotsOnTargetPercentage),
		YellowCards:             s.YellowCards,
		RedCards:                s.RedCards,
		CurrentForm:             s.CurrentForm,
		CreatedAt:               pgtypeToTime(s.CreatedAt),
		UpdatedAt:               pgtypeToTime(s.UpdatedAt),
		DeletedAt:               pgtypeToTimePtr(s.DeletedAt),
	}
}

// Helper functions for type conversions

// pgtypeToTime converts a pgtype.Timestamptz to time.Time.
// If the value is invalid, returns time.Time{} (zero value).
func pgtypeToTime(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// pgtypeToTimePtr converts a pgtype.Timestamptz to *time.Time.
// Returns nil if the value is invalid.
func pgtypeToTimePtr(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// numericToFloat64Ptr converts a pgtype.Numeric to *float64.
// Returns nil if the value is invalid or conversion fails.
func numericToFloat64Ptr(n pgtype.Numeric) *float64 {
	if !n.Valid {
		return nil
	}
	val, err := n.Float64Value()
	if err != nil {
		return nil
	}
	return &val.Float64
}
