package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type IPRateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

var ipLimiter = &IPRateLimiter{
	requests: make(map[string][]time.Time),
	limit:    100,
	window:   time.Minute,
}

func (rl *IPRateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	requests := rl.requests[ip]
	valid := make([]time.Time, 0)
	for _, t := range requests {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.limit {
		return false
	}

	rl.requests[ip] = append(valid, now)
	return true
}

func HandleRateLimitCheck(c *gin.Context) {
	ip := c.ClientIP()
	allowed := ipLimiter.Allow(ip)

	c.JSON(http.StatusOK, gin.H{
		"ip":      ip,
		"allowed": allowed,
		"limit":   ipLimiter.limit,
	})
}
