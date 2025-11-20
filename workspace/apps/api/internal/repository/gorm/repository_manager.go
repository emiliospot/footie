package gorm

import (
	"context"

	"gorm.io/gorm"

	"github.com/emiliospot/footie/api/internal/repository"
)

// GormRepositoryManager implements RepositoryManager using GORM.
type GormRepositoryManager struct {
	db *gorm.DB
	tx *gorm.DB // transaction instance

	// Repository instances
	userRepo        repository.UserRepository
	teamRepo        repository.TeamRepository
	playerRepo      repository.PlayerRepository
	matchRepo       repository.MatchRepository
	matchEventRepo  repository.MatchEventRepository
	playerStatsRepo repository.PlayerStatisticsRepository
	teamStatsRepo   repository.TeamStatisticsRepository
}

// NewRepositoryManager creates a new GORM repository manager.
func NewRepositoryManager(db *gorm.DB) repository.RepositoryManager {
	return &GormRepositoryManager{
		db: db,
	}
}

func (rm *GormRepositoryManager) getDB() *gorm.DB {
	if rm.tx != nil {
		return rm.tx
	}
	return rm.db
}

func (rm *GormRepositoryManager) User() repository.UserRepository {
	if rm.userRepo == nil {
		rm.userRepo = NewUserRepository(rm.getDB())
	}
	return rm.userRepo
}

func (rm *GormRepositoryManager) Team() repository.TeamRepository {
	if rm.teamRepo == nil {
		rm.teamRepo = NewTeamRepository(rm.getDB())
	}
	return rm.teamRepo
}

func (rm *GormRepositoryManager) Player() repository.PlayerRepository {
	if rm.playerRepo == nil {
		rm.playerRepo = NewPlayerRepository(rm.getDB())
	}
	return rm.playerRepo
}

func (rm *GormRepositoryManager) Match() repository.MatchRepository {
	if rm.matchRepo == nil {
		rm.matchRepo = NewMatchRepository(rm.getDB())
	}
	return rm.matchRepo
}

func (rm *GormRepositoryManager) MatchEvent() repository.MatchEventRepository {
	if rm.matchEventRepo == nil {
		rm.matchEventRepo = NewMatchEventRepository(rm.getDB())
	}
	return rm.matchEventRepo
}

func (rm *GormRepositoryManager) PlayerStatistics() repository.PlayerStatisticsRepository {
	if rm.playerStatsRepo == nil {
		rm.playerStatsRepo = NewPlayerStatisticsRepository(rm.getDB())
	}
	return rm.playerStatsRepo
}

func (rm *GormRepositoryManager) TeamStatistics() repository.TeamStatisticsRepository {
	if rm.teamStatsRepo == nil {
		rm.teamStatsRepo = NewTeamStatisticsRepository(rm.getDB())
	}
	return rm.teamStatsRepo
}

// BeginTx starts a new transaction.
func (rm *GormRepositoryManager) BeginTx(ctx context.Context) (repository.RepositoryManager, error) {
	tx := rm.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &GormRepositoryManager{
		db: rm.db,
		tx: tx,
	}, nil
}

// Commit commits the transaction.
func (rm *GormRepositoryManager) Commit() error {
	if rm.tx == nil {
		return gorm.ErrInvalidTransaction
	}
	return rm.tx.Commit().Error
}

// Rollback rolls back the transaction.
func (rm *GormRepositoryManager) Rollback() error {
	if rm.tx == nil {
		return gorm.ErrInvalidTransaction
	}
	return rm.tx.Rollback().Error
}
