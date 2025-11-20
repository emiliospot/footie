package gorm

import (
	"context"

	"gorm.io/gorm"

	"github.com/emiliospot/footie/api/internal/domain/models"
	"github.com/emiliospot/footie/api/internal/repository"
)

// PlayerStatistics Repository.

type GormPlayerStatisticsRepository struct {
	db *gorm.DB
}

func NewPlayerStatisticsRepository(db *gorm.DB) repository.PlayerStatisticsRepository {
	return &GormPlayerStatisticsRepository{db: db}
}

func (r *GormPlayerStatisticsRepository) Create(ctx context.Context, stats *models.PlayerStatistics) error {
	return r.db.WithContext(ctx).Create(stats).Error
}

func (r *GormPlayerStatisticsRepository) FindByID(ctx context.Context, id uint) (*models.PlayerStatistics, error) {
	var stats models.PlayerStatistics
	err := r.db.WithContext(ctx).Preload("Player").First(&stats, id).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

//nolint:dupl // Similar pattern used for Team statistics - intentional.
func (r *GormPlayerStatisticsRepository) FindByPlayerID(ctx context.Context, playerID uint, filters map[string]interface{}) ([]models.PlayerStatistics, error) {
	var stats []models.PlayerStatistics
	query := r.db.WithContext(ctx).Where("player_id = ?", playerID).Preload("Player")

	if season, ok := filters["season"].(string); ok && season != "" {
		query = query.Where("season = ?", season)
	}

	if competition, ok := filters["competition"].(string); ok && competition != "" {
		query = query.Where("competition = ?", competition)
	}

	err := query.Find(&stats).Error
	return stats, err
}

func (r *GormPlayerStatisticsRepository) Update(ctx context.Context, stats *models.PlayerStatistics) error {
	return r.db.WithContext(ctx).Save(stats).Error
}

func (r *GormPlayerStatisticsRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.PlayerStatistics{}, id).Error
}

// TeamStatistics Repository.

type GormTeamStatisticsRepository struct {
	db *gorm.DB
}

func NewTeamStatisticsRepository(db *gorm.DB) repository.TeamStatisticsRepository {
	return &GormTeamStatisticsRepository{db: db}
}

func (r *GormTeamStatisticsRepository) Create(ctx context.Context, stats *models.TeamStatistics) error {
	return r.db.WithContext(ctx).Create(stats).Error
}

func (r *GormTeamStatisticsRepository) FindByID(ctx context.Context, id uint) (*models.TeamStatistics, error) {
	var stats models.TeamStatistics
	err := r.db.WithContext(ctx).Preload("Team").First(&stats, id).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

//nolint:dupl // Similar pattern used for Player statistics - intentional.
func (r *GormTeamStatisticsRepository) FindByTeamID(ctx context.Context, teamID uint, filters map[string]interface{}) ([]models.TeamStatistics, error) {
	var stats []models.TeamStatistics
	query := r.db.WithContext(ctx).Where("team_id = ?", teamID).Preload("Team")

	if season, ok := filters["season"].(string); ok && season != "" {
		query = query.Where("season = ?", season)
	}

	if competition, ok := filters["competition"].(string); ok && competition != "" {
		query = query.Where("competition = ?", competition)
	}

	err := query.Find(&stats).Error
	return stats, err
}

func (r *GormTeamStatisticsRepository) Update(ctx context.Context, stats *models.TeamStatistics) error {
	return r.db.WithContext(ctx).Save(stats).Error
}

func (r *GormTeamStatisticsRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.TeamStatistics{}, id).Error
}
