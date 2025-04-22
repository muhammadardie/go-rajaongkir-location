package middleware

import (
	"fmt"
	"go-rajaongkir-location/utils/response"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	requests   map[string]int
	lastReset  map[string]time.Time
	limit      int
	windowSize time.Duration
	mu         sync.Mutex
}

// NewRateLimiterFromEnv creates a new rate limiter with values from environment variables
func CreateRateLimiter() *RateLimiter {
	// Get rate limit values from environment with defaults
	limitStr := os.Getenv("RATE_REQUEST")
	minutesStr := os.Getenv("RATE_MINUTE")

	// Set defaults if environment variables are not set
	limit := 10  // Default limit
	minutes := 1 // Default window size in minutes

	// Parse environment variables
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if minutesStr != "" {
		if parsedMinutes, err := strconv.Atoi(minutesStr); err == nil && parsedMinutes > 0 {
			minutes = parsedMinutes
		}
	}

	windowSize := time.Duration(minutes) * time.Minute

	return &RateLimiter{
		requests:   make(map[string]int),
		lastReset:  make(map[string]time.Time),
		limit:      limit,
		windowSize: windowSize,
	}
}

// Keep the original constructor for flexibility
func NewRateLimiter(limit int, windowSize time.Duration) *RateLimiter {
	return &RateLimiter{
		requests:   make(map[string]int),
		lastReset:  make(map[string]time.Time),
		limit:      limit,
		windowSize: windowSize,
	}
}

func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			apiKey = c.ClientIP() // Fallback to IP address
		}

		rl.mu.Lock()
		defer rl.mu.Unlock()

		now := time.Now()
		lastResetTime, exists := rl.lastReset[apiKey]

		// Calculate time remaining until reset
		var timeRemaining time.Duration
		if exists {
			timeRemaining = rl.windowSize - now.Sub(lastResetTime)
			if timeRemaining < 0 {
				timeRemaining = 0
			}
		}

		if !exists || now.Sub(lastResetTime) > rl.windowSize {
			// Reset counter for new window
			rl.requests[apiKey] = 1
			rl.lastReset[apiKey] = now
			timeRemaining = rl.windowSize
		} else {
			// Increment counter
			rl.requests[apiKey]++
		}

		// Add headers to show rate limit status
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", max(0, rl.limit-rl.requests[apiKey])))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%.0f", timeRemaining.Seconds()))

		// Check if limit exceeded
		if rl.requests[apiKey] > rl.limit {
			message := fmt.Sprintf("Rate limit exceeded, retry after %.0f seconds", timeRemaining.Seconds())
			response.ErrorResponse(c, message, http.StatusTooManyRequests)
			c.Abort()
			return
		}

		c.Next()
	}
}
