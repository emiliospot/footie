package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	redisClient "github.com/redis/go-redis/v9"

	"github.com/emiliospot/footie/api/internal/api"
	"github.com/emiliospot/footie/api/internal/config"
	"github.com/emiliospot/footie/api/internal/infrastructure/database"
	"github.com/emiliospot/footie/api/internal/infrastructure/logger"
	"github.com/emiliospot/footie/api/internal/infrastructure/redis"
	ws "github.com/emiliospot/footie/api/internal/infrastructure/websocket"
)

// @title Footie API.
// @version 1.0.
// @description Football Analytics Platform API.
// @termsOfService http://swagger.io/terms/.

// @contact.name API Support.
// @contact.url http://www.footie.com/support.
// @contact.email support@footie.com.

// @license.name MIT.
// @license.url https://opensource.org/licenses/MIT.

// @host localhost:8080.
// @BasePath /api/v1.

// @securityDefinitions.apikey BearerAuth.
// @in header.
// @name Authorization.
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	appLogger := logger.NewLogger(cfg.Log.Level, cfg.Log.Format)
	appLogger.Info("Starting Footie API", "version", cfg.App.Version, "environment", cfg.App.Environment)

	// Initialize context
	ctx := context.Background()

	// Database connection (optional in development for mock data endpoints)
	var pool *pgxpool.Pool
	if cfg.IsDevelopment() && os.Getenv("SKIP_DB") == "true" {
		appLogger.Warn("Skipping database connection (SKIP_DB=true). Mock data endpoints will work, but database-dependent endpoints will fail.")
	} else {
		// Run database migrations first
		appLogger.Info("Running database migrations...")
		databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
			cfg.Database.SSLMode,
		)
		migrationsPath := "./migrations" // Relative to apps/api directory
		if migErr := database.RunMigrations(databaseURL, migrationsPath); migErr != nil {
			if cfg.IsDevelopment() {
				appLogger.Warn("Failed to run migrations (database may not be running)", "error", migErr)
				appLogger.Warn("To skip database, set SKIP_DB=true. To start database, run: docker-compose -f workspace/infra/docker/docker-compose.yml up -d postgres redis")
			} else {
				appLogger.Fatal("Failed to run migrations", "error", migErr)
			}
		} else {
			appLogger.Info("Database migrations completed successfully")
		}

		// Initialize pgx connection pool
		pgxCfg := &database.PgxConfig{
			Host:     cfg.Database.Host,
			Port:     cfg.Database.Port,
			User:     cfg.Database.User,
			Password: cfg.Database.Password,
			Database: cfg.Database.Name,
			SSLMode:  cfg.Database.SSLMode,
		}

		var dbErr error
		pool, dbErr = database.NewPgxPool(ctx, pgxCfg)
		if dbErr != nil {
			if cfg.IsDevelopment() {
				appLogger.Warn("Failed to connect to database (database may not be running)", "error", dbErr)
				appLogger.Warn("To skip database, set SKIP_DB=true. To start database, run: docker-compose -f workspace/infra/docker/docker-compose.yml up -d postgres redis")
			} else {
				appLogger.Fatal("Failed to connect to database", "error", dbErr)
			}
		} else {
			defer pool.Close()
			appLogger.Info("Database connected successfully", "max_conns", pool.Config().MaxConns)
		}
	}

	// Initialize Redis (optional in development)
	var redisClient *redisClient.Client
	if cfg.IsDevelopment() && os.Getenv("SKIP_REDIS") == "true" {
		appLogger.Warn("Skipping Redis connection (SKIP_REDIS=true). Real-time features will not work.")
	} else {
		var redisErr error
		redisClient, redisErr = redis.NewRedisClient(cfg.Redis)
		if redisErr != nil {
			if cfg.IsDevelopment() {
				appLogger.Warn("Failed to connect to Redis (Redis may not be running)", "error", redisErr)
				appLogger.Warn("To skip Redis, set SKIP_REDIS=true. To start Redis, run: docker-compose -f workspace/infra/docker/docker-compose.yml up -d redis")
			} else {
				appLogger.Fatal("Failed to connect to Redis", "error", redisErr)
			}
		} else {
			appLogger.Info("Redis connected successfully")
		}
	}

	// Initialize WebSocket hub (only if Redis is available)
	var hub *ws.Hub
	if redisClient != nil {
		hub = ws.NewHub(redisClient, appLogger)
		go hub.Run(ctx)
		appLogger.Info("WebSocket hub started")
	} else {
		appLogger.Warn("WebSocket hub not started (Redis not available)")
		hub = nil
	}

	// Initialize router (pool and redis can be nil in development for mock endpoints)
	// Note: Handlers that use database will fail if pool is nil, but rankings (mock data) will work
	router := api.NewRouter(cfg, pool, redisClient, hub, appLogger)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.API.Host, cfg.API.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		appLogger.Info("Server starting", "address", srv.Addr)
		if srvErr := srv.ListenAndServe(); srvErr != nil && srvErr != http.ErrServerClosed {
			appLogger.Fatal("Failed to start server", "error", srvErr)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if shutdownErr := srv.Shutdown(ctx); shutdownErr != nil {
		appLogger.Fatal("Server forced to shutdown", "error", shutdownErr)
	}

	// Close database connection pool (if connected)
	if pool != nil {
		pool.Close()
	}

	// Close Redis connection (if connected)
	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			appLogger.Error("Error closing Redis connection", "error", err)
		}
	}

	appLogger.Info("Server exited successfully")
}
