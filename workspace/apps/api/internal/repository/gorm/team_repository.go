package gorm

import (
	"context"

	"gorm.io/gorm"

	"github.com/emiliospot/footie/api/internal/domain/models"
	"github.com/emiliospot/footie/api/internal/repository"
)

type GormTeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) repository.TeamRepository {
	return &GormTeamRepository{db: db}
}

func (r *GormTeamRepository) Create(ctx context.Context, team *models.Team) error {
	return r.db.WithContext(ctx).Create(team).Error
}

func (r *GormTeamRepository) FindByID(ctx context.Context, id uint) (*models.Team, error) {
	var team models.Team
	err := r.db.WithContext(ctx).Preload("Players").First(&team, id).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *GormTeamRepository) FindByCode(ctx context.Context, code string) (*models.Team, error) {
	var team models.Team
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&team).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *GormTeamRepository) Update(ctx context.Context, team *models.Team) error {
	return r.db.WithContext(ctx).Save(team).Error
}

func (r *GormTeamRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Team{}, id).Error
}

func (r *GormTeamRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]models.Team, int64, error) {
	var teams []models.Team
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Team{})

	// Apply filters
	if country, ok := filters["country"].(string); ok && country != "" {
		query = query.Where("country = ?", country)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&teams).Error; err != nil {
		return nil, 0, err
	}

	return teams, total, nil
}
