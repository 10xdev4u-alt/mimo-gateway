package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ResponseCache struct {
	mu      sync.RWMutex
	entries map[string]CacheItem
}

type CacheItem struct {
	Response  []byte
	ExpiresAt time.Time
}

var responseCache = &ResponseCache{entries: make(map[string]CacheItem)}

func CacheMiddleware(ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		key := c.Request.URL.String()
		responseCache.mu.RLock()
		item, found := responseCache.entries[key]
		responseCache.mu.RUnlock()

		if found && time.Now().Before(item.ExpiresAt) {
			c.Data(http.StatusOK, "application/json", item.Response)
			c.Abort()
			return
		}

		c.Next()

		if c.Writer.Status() == http.StatusOK {
			body := c.Writer.Body.Bytes()
			responseCache.mu.Lock()
			responseCache.entries[key] = CacheItem{
				Response:  body,
				ExpiresAt: time.Now().Add(ttl),
			}
			responseCache.mu.Unlock()
		}
	}
}
