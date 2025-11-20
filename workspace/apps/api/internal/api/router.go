package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// "github.com/emiliospot/footie/api/internal/api/handlers" // TODO: Update for sqlc
	"github.com/emiliospot/footie/api/internal/api/middleware"
	"github.com/emiliospot/footie/api/internal/config"
	"github.com/emiliospot/footie/api/internal/infrastructure/logger"
)

// NewRouter creates and configures the HTTP router.
func NewRouter(cfg *config.Config, pool *pgxpool.Pool, redis *redis.Client, logger *logger.Logger) *gin.Engine {
	// Set Gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.RequestID())

	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           12 * 3600, // 12 hours
	}
	router.Use(cors.New(corsConfig))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"version": cfg.App.Version,
		})
	})

	// Swagger documentation
	if cfg.IsDevelopment() {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// TODO: Update handlers to use sqlc queries instead of GORM
	// The handlers need to be refactored to work with sqlc + pgx
	// For now, we'll just set up the basic routes structure
	
	// Initialize sqlc queries
	// queries := sqlc.New(pool)
	
	// Temporarily commented out until handlers are updated for sqlc
	/*
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(cfg, pool, logger)
	userHandler := handlers.NewUserHandler(pool, logger)
	teamHandler := handlers.NewTeamHandler(pool, logger)
	playerHandler := handlers.NewPlayerHandler(pool, logger)
	matchHandler := handlers.NewMatchHandler(pool, logger)

	// API v1 routes
	v1 := router.Group("/api/v1")
	// Public routes
	auth := v1.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.RefreshToken)

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))

	// User routes
	users := protected.Group("/users")
	users.GET("/me", userHandler.GetCurrentUser)
	users.PUT("/me", userHandler.UpdateCurrentUser)
	users.GET("/:id", userHandler.GetUser)

	// Team routes
	teams := protected.Group("/teams")
	teams.GET("", teamHandler.ListTeams)
	teams.GET("/:id", teamHandler.GetTeam)
	teams.POST("", middleware.RequireRole("analyst"), teamHandler.CreateTeam)
	teams.PUT("/:id", middleware.RequireRole("analyst"), teamHandler.UpdateTeam)
	teams.DELETE("/:id", middleware.RequireRole("admin"), teamHandler.DeleteTeam)
	teams.GET("/:id/players", teamHandler.GetTeamPlayers)
	teams.GET("/:id/statistics", teamHandler.GetTeamStatistics)

	// Player routes
	players := protected.Group("/players")
	players.GET("", playerHandler.ListPlayers)
	players.GET("/:id", playerHandler.GetPlayer)
	players.POST("", middleware.RequireRole("analyst"), playerHandler.CreatePlayer)
	players.PUT("/:id", middleware.RequireRole("analyst"), playerHandler.UpdatePlayer)
	players.DELETE("/:id", middleware.RequireRole("admin"), playerHandler.DeletePlayer)
	players.GET("/:id/statistics", playerHandler.GetPlayerStatistics)

	// Match routes
	matches := protected.Group("/matches")
	matches.GET("", matchHandler.ListMatches)
	matches.GET("/:id", matchHandler.GetMatch)
	matches.POST("", middleware.RequireRole("analyst"), matchHandler.CreateMatch)
	matches.PUT("/:id", middleware.RequireRole("analyst"), matchHandler.UpdateMatch)
	matches.DELETE("/:id", middleware.RequireRole("admin"), matchHandler.DeleteMatch)
	matches.GET("/:id/events", matchHandler.GetMatchEvents)
	matches.POST("/:id/events", middleware.RequireRole("analyst"), matchHandler.CreateMatchEvent)

	// Admin routes
	admin := protected.Group("/admin")
	admin.Use(middleware.RequireRole("admin"))
	admin.GET("/users", userHandler.ListUsers)
	admin.PUT("/users/:id/role", userHandler.UpdateUserRole)
	admin.DELETE("/users/:id", userHandler.DeleteUser)
	*/

	return router
}
