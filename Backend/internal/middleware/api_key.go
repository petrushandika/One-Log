package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/petrushandika/one-log/internal/repository"
	"github.com/petrushandika/one-log/pkg/utils"
)

// APIKeyAuth validates the X-API-Key header against registered sources
func APIKeyAuth(repo repository.SourceRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		
		// Fallback checking Authorization header just in case "Bearer <token>" is passed
		if apiKey == "" {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				apiKey = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if apiKey == "" {
			utils.Error(c, http.StatusUnauthorized, "Missing API Key", nil)
			c.Abort()
			return
		}

		// Hash incoming API key to match it with DB stored format
		hashedApiKey := utils.HashAPIKey(apiKey)

		// Find in DB
		source, err := repo.FindByAPIKey(hashedApiKey)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "Failed to validate API Key", err.Error())
			c.Abort()
			return
		}

		if source == nil {
			utils.Error(c, http.StatusUnauthorized, "Invalid API Key", nil)
			c.Abort()
			return
		}

		// Inject Source information into Context (so the handler knows who owns this log)
		c.Set("source_id", source.ID)
		
		c.Next()
	}
}
