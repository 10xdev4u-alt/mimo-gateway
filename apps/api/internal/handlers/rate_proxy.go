package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ProxyRateLimiter struct {
	mu       sync.RWMutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

var proxyLimiter = &ProxyRateLimiter{
	requests: make(map[string][]time.Time),
	limit:    100,
	window:   time.Minute,
}

func HandleProxyRateLimit(c *gin.Context) {
	ip := c.ClientIP()
	proxyLimiter.mu.RLock()
	count := len(proxyLimiter.requests[ip])
	proxyLimiter.mu.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"ip":        ip,
		"count":     count,
		"limit":     proxyLimiter.limit,
		"window":    proxyLimiter.window.String(),
		"remaining": proxyLimiter.limit - count,
	})
}
