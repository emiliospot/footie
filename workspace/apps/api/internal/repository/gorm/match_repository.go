package gorm

import (
	"context"

	"gorm.io/gorm"

	"github.com/emiliospot/footie/api/internal/domain/models"
	"github.com/emiliospot/footie/api/internal/repository"
)

type GormMatchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) repository.MatchRepository {
	return &GormMatchRepository{db: db}
}

func (r *GormMatchRepository) Create(ctx context.Context, match *models.Match) error {
	return r.db.WithContext(ctx).Create(match).Error
}

func (r *GormMatchRepository) FindByID(ctx context.Context, id uint) (*models.Match, error) {
	var match models.Match
	err := r.db.WithContext(ctx).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("Events").
		First(&match, id).Error
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *GormMatchRepository) Update(ctx context.Context, match *models.Match) error {
	return r.db.WithContext(ctx).Save(match).Error
}

func (r *GormMatchRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Match{}, id).Error
}

func (r *GormMatchRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]models.Match, int64, error) {
	var matches []models.Match
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Match{}).
		Preload("HomeTeam").
		Preload("AwayTeam")

	if teamID, ok := filters["team_id"].(uint); ok && teamID > 0 {
		query = query.Where("home_team_id = ? OR away_team_id = ?", teamID, teamID)
	}

	if season, ok := filters["season"].(string); ok && season != "" {
		query = query.Where("season = ?", season)
	}

	if competition, ok := filters["competition"].(string); ok && competition != "" {
		query = query.Where("competition = ?", competition)
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("match_date DESC").Offset(offset).Limit(limit).Find(&matches).Error; err != nil {
		return nil, 0, err
	}

	return matches, total, nil
}
