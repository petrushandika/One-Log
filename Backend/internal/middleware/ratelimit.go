package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/petrushandika/one-log/pkg/utils"
)

// RateLimiterConfig configures the rate limiter
type RateLimiterConfig struct {
	// Requests per window
	Requests int
	// Time window duration
	Window time.Duration
	// Key function to extract identifier from request (e.g., IP, API Key, User ID)
	KeyFunc func(*gin.Context) string
}

// RateLimiter implements token bucket algorithm for rate limiting
type RateLimiter struct {
	config  RateLimiterConfig
	buckets map[string]*bucket
	mu      sync.RWMutex
}

type bucket struct {
	tokens    int
	lastCheck time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	if config.Requests <= 0 {
		config.Requests = 100
	}
	if config.Window <= 0 {
		config.Window = time.Minute
	}

	rl := &RateLimiter{
		config:  config,
		buckets: make(map[string]*bucket),
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// Allow checks if a request is allowed
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, exists := rl.buckets[key]

	if !exists {
		// Create new bucket
		rl.buckets[key] = &bucket{
			tokens:    rl.config.Requests - 1,
			lastCheck: now,
		}
		return true
	}

	// Calculate tokens to add based on time passed
	elapsed := now.Sub(b.lastCheck)
	tokensToAdd := int(elapsed/rl.config.Window) * rl.config.Requests

	if tokensToAdd > 0 {
		b.tokens = min(b.tokens+tokensToAdd, rl.config.Requests)
		b.lastCheck = now
	}

	if b.tokens > 0 {
		b.tokens--
		return true
	}

	return false
}

// GetRetryAfter returns how long to wait before next request
func (rl *RateLimiter) GetRetryAfter(key string) time.Duration {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	b, exists := rl.buckets[key]
	if !exists {
		return 0
	}

	if b.tokens > 0 {
		return 0
	}

	// Calculate time until next token
	elapsed := time.Since(b.lastCheck)
	timePerToken := rl.config.Window / time.Duration(rl.config.Requests)
	return timePerToken - elapsed%timePerToken
}

// cleanup removes old buckets periodically
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, b := range rl.buckets {
			if now.Sub(b.lastCheck) > rl.config.Window*2 {
				delete(rl.buckets, key)
			}
		}
		rl.mu.Unlock()
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RateLimitMiddleware creates a Gin middleware for rate limiting
func RateLimitMiddleware(config RateLimiterConfig) gin.HandlerFunc {
	limiter := NewRateLimiter(config)

	return func(c *gin.Context) {
		key := config.KeyFunc(c)
		if key == "" {
			key = c.ClientIP()
		}

		if !limiter.Allow(key) {
			retryAfter := limiter.GetRetryAfter(key)
			c.Header("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", config.Requests))
			c.Header("X-RateLimit-Window", config.Window.String())

			utils.Error(c, http.StatusTooManyRequests, "Rate limit exceeded. Please try again later.", map[string]interface{}{
				"retry_after_seconds": int(retryAfter.Seconds()),
				"limit":               config.Requests,
				"window":              config.Window.String(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByAPIKey creates rate limiter keyed by API Key (for ingestion endpoint)
func RateLimitByAPIKey(requests int, window time.Duration) gin.HandlerFunc {
	return RateLimitMiddleware(RateLimiterConfig{
		Requests: requests,
		Window:   window,
		KeyFunc: func(c *gin.Context) string {
			// Get API key from header
			apiKey := c.GetHeader("X-API-Key")
			if apiKey == "" {
				authHeader := c.GetHeader("Authorization")
				if strings.HasPrefix(authHeader, "Bearer ") {
					apiKey = strings.TrimPrefix(authHeader, "Bearer ")
				}
			}
			return apiKey
		},
	})
}

// RateLimitByIP creates rate limiter keyed by IP address (for public endpoints)
func RateLimitByIP(requests int, window time.Duration) gin.HandlerFunc {
	return RateLimitMiddleware(RateLimiterConfig{
		Requests: requests,
		Window:   window,
		KeyFunc: func(c *gin.Context) string {
			return c.ClientIP()
		},
	})
}

// RateLimitByJWT creates rate limiter keyed by JWT user ID (for authenticated endpoints)
func RateLimitByJWT(requests int, window time.Duration) gin.HandlerFunc {
	return RateLimitMiddleware(RateLimiterConfig{
		Requests: requests,
		Window:   window,
		KeyFunc: func(c *gin.Context) string {
			userID, exists := c.Get("user_id")
			if !exists {
				return c.ClientIP() // Fallback to IP
			}
			return fmt.Sprintf("user_%v", userID)
		},
	})
}
