package gorm

import (
	"context"

	"gorm.io/gorm"

	"github.com/emiliospot/footie/api/internal/domain/models"
	"github.com/emiliospot/footie/api/internal/repository"
)

type GormMatchEventRepository struct {
	db *gorm.DB
}

func NewMatchEventRepository(db *gorm.DB) repository.MatchEventRepository {
	return &GormMatchEventRepository{db: db}
}

func (r *GormMatchEventRepository) Create(ctx context.Context, event *models.MatchEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *GormMatchEventRepository) FindByID(ctx context.Context, id uint) (*models.MatchEvent, error) {
	var event models.MatchEvent
	err := r.db.WithContext(ctx).
		Preload("Player").
		Preload("Team").
		Preload("SecondaryPlayer").
		First(&event, id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *GormMatchEventRepository) FindByMatchID(ctx context.Context, matchID uint) ([]models.MatchEvent, error) {
	var events []models.MatchEvent
	err := r.db.WithContext(ctx).
		Where("match_id = ?", matchID).
		Preload("Player").
		Preload("Team").
		Preload("SecondaryPlayer").
		Order("minute ASC, extra_minute ASC").
		Find(&events).Error
	return events, err
}

func (r *GormMatchEventRepository) Update(ctx context.Context, event *models.MatchEvent) error {
	return r.db.WithContext(ctx).Save(event).Error
}

func (r *GormMatchEventRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.MatchEvent{}, id).Error
}
