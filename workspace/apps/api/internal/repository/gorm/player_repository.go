package gorm

import (
	"context"

	"gorm.io/gorm"

	"github.com/emiliospot/footie/api/internal/domain/models"
	"github.com/emiliospot/footie/api/internal/repository"
)

type GormPlayerRepository struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) repository.PlayerRepository {
	return &GormPlayerRepository{db: db}
}

func (r *GormPlayerRepository) Create(ctx context.Context, player *models.Player) error {
	return r.db.WithContext(ctx).Create(player).Error
}

func (r *GormPlayerRepository) FindByID(ctx context.Context, id uint) (*models.Player, error) {
	var player models.Player
	err := r.db.WithContext(ctx).Preload("Team").First(&player, id).Error
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *GormPlayerRepository) Update(ctx context.Context, player *models.Player) error {
	return r.db.WithContext(ctx).Save(player).Error
}

func (r *GormPlayerRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Player{}, id).Error
}

func (r *GormPlayerRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]models.Player, int64, error) {
	var players []models.Player
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Player{}).Preload("Team")

	if teamID, ok := filters["team_id"].(uint); ok && teamID > 0 {
		query = query.Where("team_id = ?", teamID)
	}

	if position, ok := filters["position"].(string); ok && position != "" {
		query = query.Where("position = ?", position)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&players).Error; err != nil {
		return nil, 0, err
	}

	return players, total, nil
}

func (r *GormPlayerRepository) FindByTeamID(ctx context.Context, teamID uint) ([]models.Player, error) {
	var players []models.Player
	err := r.db.WithContext(ctx).Where("team_id = ?", teamID).Find(&players).Error
	return players, err
}
