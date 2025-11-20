package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PgxConfig holds configuration for pgx connection pool.
type PgxConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

// NewPgxPool creates a new pgx connection pool.
func NewPgxPool(ctx context.Context, cfg *PgxConfig) (*pgxpool.Pool, error) {
	// Build connection string
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.SSLMode,
	)

	// Parse config
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %w", err)
	}

	// Configure pool settings for production workloads
	poolConfig.MaxConns = 25                               // Maximum number of connections
	poolConfig.MinConns = 5                                // Minimum number of connections
	poolConfig.MaxConnLifetime = time.Hour                 // Max connection lifetime
	poolConfig.MaxConnIdleTime = 30 * time.Minute          // Max idle time
	poolConfig.HealthCheckPeriod = time.Minute             // Health check interval
	poolConfig.ConnConfig.ConnectTimeout = 5 * time.Second // Connection timeout

	// Create pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Ping to verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}
