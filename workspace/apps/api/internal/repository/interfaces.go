package repository

import (
	"context"

	"github.com/emiliospot/footie/api/internal/domain/models"
)

// Database abstraction layer - allows switching ORMs easily.

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]models.User, int64, error)
}

// TeamRepository defines the interface for team data operations.
type TeamRepository interface {
	Create(ctx context.Context, team *models.Team) error
	FindByID(ctx context.Context, id uint) (*models.Team, error)
	FindByCode(ctx context.Context, code string) (*models.Team, error)
	Update(ctx context.Context, team *models.Team) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]models.Team, int64, error)
}

// PlayerRepository defines the interface for player data operations.
type PlayerRepository interface {
	Create(ctx context.Context, player *models.Player) error
	FindByID(ctx context.Context, id uint) (*models.Player, error)
	Update(ctx context.Context, player *models.Player) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]models.Player, int64, error)
	FindByTeamID(ctx context.Context, teamID uint) ([]models.Player, error)
}

// MatchRepository defines the interface for match data operations.
type MatchRepository interface {
	Create(ctx context.Context, match *models.Match) error
	FindByID(ctx context.Context, id uint) (*models.Match, error)
	Update(ctx context.Context, match *models.Match) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]models.Match, int64, error)
}

// MatchEventRepository defines the interface for match event data operations.
type MatchEventRepository interface {
	Create(ctx context.Context, event *models.MatchEvent) error
	FindByID(ctx context.Context, id uint) (*models.MatchEvent, error)
	FindByMatchID(ctx context.Context, matchID uint) ([]models.MatchEvent, error)
	Update(ctx context.Context, event *models.MatchEvent) error
	Delete(ctx context.Context, id uint) error
}

// PlayerStatisticsRepository defines the interface for player statistics operations.
//
//nolint:dupl // Similar interface pattern for Team statistics - intentional.
type PlayerStatisticsRepository interface {
	Create(ctx context.Context, stats *models.PlayerStatistics) error
	FindByID(ctx context.Context, id uint) (*models.PlayerStatistics, error)
	FindByPlayerID(ctx context.Context, playerID uint, filters map[string]interface{}) ([]models.PlayerStatistics, error)
	Update(ctx context.Context, stats *models.PlayerStatistics) error
	Delete(ctx context.Context, id uint) error
}

// TeamStatisticsRepository defines the interface for team statistics operations.
//
//nolint:dupl // Similar interface pattern for Player statistics - intentional.
type TeamStatisticsRepository interface {
	Create(ctx context.Context, stats *models.TeamStatistics) error
	FindByID(ctx context.Context, id uint) (*models.TeamStatistics, error)
	FindByTeamID(ctx context.Context, teamID uint, filters map[string]interface{}) ([]models.TeamStatistics, error)
	Update(ctx context.Context, stats *models.TeamStatistics) error
	Delete(ctx context.Context, id uint) error
}

// RepositoryManager provides access to all repositories.
type RepositoryManager interface {
	User() UserRepository
	Team() TeamRepository
	Player() PlayerRepository
	Match() MatchRepository
	MatchEvent() MatchEventRepository
	PlayerStatistics() PlayerStatisticsRepository
	TeamStatistics() TeamStatisticsRepository

	// Transaction support
	BeginTx(ctx context.Context) (RepositoryManager, error)
	Commit() error
	Rollback() error
}
