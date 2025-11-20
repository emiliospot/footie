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

	"github.com/emiliospot/footie/api/internal/api"
	"github.com/emiliospot/footie/api/internal/config"
	"github.com/emiliospot/footie/api/internal/infrastructure/database"
	"github.com/emiliospot/footie/api/internal/infrastructure/logger"
	"github.com/emiliospot/footie/api/internal/infrastructure/redis"
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

	// Run database migrations first
	appLogger.Info("Running database migrations...")
	if migErr := database.RunMigrations(cfg.Database.URL, "migrations"); migErr != nil {
		appLogger.Fatal("Failed to run migrations", "error", migErr)
	}
	appLogger.Info("Database migrations completed successfully")

	// Initialize pgx connection pool
	ctx := context.Background()
	pgxCfg := &database.PgxConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Database: cfg.Database.Name,
		SSLMode:  cfg.Database.SSLMode,
	}
	
	pool, err := database.NewPgxPool(ctx, pgxCfg)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", "error", err)
	}
	defer pool.Close()
	appLogger.Info("Database connected successfully", "max_conns", pool.Config().MaxConns)

	// Initialize Redis
	redisClient, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		appLogger.Fatal("Failed to connect to Redis", "error", err)
	}
	appLogger.Info("Redis connected successfully")

	// Initialize router
	router := api.NewRouter(cfg, pool, redisClient, appLogger)

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

	// Close database connection pool
	pool.Close()

	// Close Redis connection
	if err := redisClient.Close(); err != nil {
		appLogger.Error("Error closing Redis connection", "error", err)
	}

	appLogger.Info("Server exited successfully")
}
