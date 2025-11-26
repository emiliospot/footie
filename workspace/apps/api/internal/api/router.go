package api

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/emiliospot/footie/api/internal/api/handlers"
	"github.com/emiliospot/footie/api/internal/api/middleware"
	"github.com/emiliospot/footie/api/internal/config"
	"github.com/emiliospot/footie/api/internal/infrastructure/logger"
	"github.com/emiliospot/footie/api/internal/infrastructure/webhooks"
	"github.com/emiliospot/footie/api/internal/infrastructure/webhooks/providers"
	ws "github.com/emiliospot/footie/api/internal/infrastructure/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, check against allowed origins
		return true
	},
}

// NewRouter creates and configures the HTTP router.
func NewRouter(cfg *config.Config, pool *pgxpool.Pool, redis *redis.Client, hub *ws.Hub, logger *logger.Logger) *gin.Engine {
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

	// Initialize base handler with common dependencies
	baseHandler := handlers.NewBaseHandler(cfg, pool, redis, logger)

	// Initialize webhook provider registry
	providerRegistry := webhooks.NewRegistry()
	providerRegistry.Register(providers.NewGenericProvider())
	providerRegistry.Register(providers.NewOptaProvider())
	providerRegistry.Register(providers.NewStatsBombProvider())

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(baseHandler)
	matchHandler := handlers.NewMatchHandler(baseHandler)
	rankingsHandler := handlers.NewRankingsHandler(baseHandler)
	webhookHandler := handlers.NewWebhookHandler(baseHandler, &cfg.Webhook, providerRegistry)

	// Health check endpoint
	router.GET("/health", healthHandler.Check)

	// Webhook endpoints (public, but signature-verified)
	webhooks := router.Group("/webhooks")
	webhooks.POST("/matches", webhookHandler.HandleMatchEvents)
	webhooks.POST("/matches/:id/status", webhookHandler.HandleMatchStatus)

	// WebSocket endpoint for real-time match updates
	router.GET("/ws/matches/:id", func(c *gin.Context) {
		matchIDStr := c.Param("id")
		matchID, err := strconv.ParseInt(matchIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
			return
		}

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logger.Error("Failed to upgrade WebSocket", "error", err)
			return
		}

		// Get user ID from context (if authenticated)
		userID := int32(0)
		if userIDVal, exists := c.Get("user_id"); exists {
			if uid, ok := userIDVal.(int32); ok {
				userID = uid
			}
		}

		// Serve WebSocket connection
		ws.ServeWs(hub, conn, int32(matchID), userID)
	})

	// Swagger documentation
	if cfg.IsDevelopment() {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// API v1 routes
	v1 := router.Group("/api/v1")

	// Public routes (no authentication required)
	// TODO: Implement auth handler (register, login, refresh)
	// auth := v1.Group("/auth")
	// auth.POST("/register", authHandler.Register)
	// auth.POST("/login", authHandler.Login)
	// auth.POST("/refresh", authHandler.RefreshToken)

	// Protected routes (authentication required)
	// For now, we'll make match routes public for development
	// In production, add: protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	protected := v1.Group("")
	// protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))

	// Match routes
	matches := protected.Group("/matches")
	matches.GET("", matchHandler.ListMatches)
	matches.GET("/:id", matchHandler.GetMatch)
	matches.GET("/:id/events", matchHandler.GetMatchEvents)
	matches.POST("/:id/events", matchHandler.CreateMatchEvent) // TODO: Add RequireRole("analyst")

	// Rankings routes
	rankings := protected.Group("/rankings")
	rankings.GET("", rankingsHandler.GetCompetitionRankings)

	// TODO: Implement additional handlers
	// - User handler (users CRUD, profile management)
	// - Team handler (teams CRUD, statistics)
	// - Player handler (players CRUD, statistics)
	// - Auth handler (JWT authentication)
	// - Admin routes (user management)

	return router
}
