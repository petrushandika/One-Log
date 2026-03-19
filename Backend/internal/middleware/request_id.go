package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDMiddleware adds a unique request ID to each request for tracing
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate new UUID for request ID
		requestID := uuid.New().String()

		// Set in context for use in handlers
		c.Set("request_id", requestID)

		// Set in response header for client-side tracing
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}
