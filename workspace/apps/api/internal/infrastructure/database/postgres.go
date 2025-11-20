package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/emiliospot/footie/api/internal/config"
	"github.com/emiliospot/footie/api/internal/domain/models"
)

// NewPostgresDB creates a new PostgreSQL database connection.
func NewPostgresDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.URL

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// RunMigrations runs database migrations.
func RunMigrations(db *gorm.DB) error {
	// Auto migrate models
	return db.AutoMigrate(
		&models.User{},
		&models.Team{},
		&models.Player{},
		&models.Match{},
		&models.MatchEvent{},
		&models.PlayerStatistics{},
		&models.TeamStatistics{},
	)
}

// Close closes the database connection.
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
