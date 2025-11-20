package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/emiliospot/footie/api/internal/infrastructure/logger"
)

// Logger creates a logging middleware.
func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request
		duration := time.Since(start)
		status := c.Writer.Status()

		log.Info("HTTP Request",
			"method", c.Request.Method,
			"path", path,
			"query", query,
			"status", status,
			"duration_ms", duration.Milliseconds(),
			"client_ip", c.ClientIP(),
			"request_id", c.GetString("request_id"),
		)

		// Log errors if any
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Error("Request error", "error", err.Error())
			}
		}
	}
}
