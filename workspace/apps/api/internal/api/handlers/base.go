package handlers

import (
	"github.com/emiliospot/footie/api/internal/config"
	"github.com/emiliospot/footie/api/internal/infrastructure/events"
	"github.com/emiliospot/footie/api/internal/infrastructure/logger"
	"github.com/emiliospot/footie/api/internal/repository/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// BaseHandler contains common dependencies for all handlers.
type BaseHandler struct {
	cfg       *config.Config
	pool      *pgxpool.Pool
	queries   *sqlc.Queries
	redis     *redis.Client
	publisher *events.Publisher
	logger    *logger.Logger
}

// NewBaseHandler creates a new base handler with common dependencies.
func NewBaseHandler(
	cfg *config.Config,
	pool *pgxpool.Pool,
	redis *redis.Client,
	logger *logger.Logger,
) *BaseHandler {
	queries := sqlc.New(pool)
	publisher := events.NewPublisher(redis, logger)

	return &BaseHandler{
		cfg:       cfg,
		pool:      pool,
		queries:   queries,
		redis:     redis,
		publisher: publisher,
		logger:    logger,
	}
}
