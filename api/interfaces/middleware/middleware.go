package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thalesmacedo1/freterapido-backend-api/infrastructure/logger"
)

func LoggerMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()

		// Log format: Method, Path, Status, Latency, ClientIP
		log.Infof("Method=%s, Path=%s, Status=%d, Latency=%v, ClientIP=%s",
			c.Request.Method,
			c.Request.URL.Path,
			statusCode,
			latency,
			c.ClientIP(),
		)
	}
}
