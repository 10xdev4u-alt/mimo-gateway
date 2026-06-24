package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	mu       sync.RWMutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

var limiter = &RateLimiter{
	requests: make(map[string][]time.Time),
	limit:    100,
	window:   time.Minute,
}

func HandleRateLimitStatus(c *gin.Context) {
	ip := c.ClientIP()
	limiter.mu.RLock()
	count := len(limiter.requests[ip])
	limiter.mu.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"ip":       ip,
		"count":    count,
		"limit":    limiter.limit,
		"window":   limiter.window.String(),
		"remaining": limiter.limit - count,
	})
}
