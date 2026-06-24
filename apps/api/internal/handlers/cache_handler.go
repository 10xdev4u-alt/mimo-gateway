package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type CacheEntry struct {
	Value     interface{}
	ExpiresAt time.Time
}

type SimpleCache struct {
	mu      sync.RWMutex
	entries map[string]CacheEntry
}

var cache = &SimpleCache{entries: make(map[string]CacheEntry)}

func (c *SimpleCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	if !ok || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	return entry.Value, true
}

func (c *SimpleCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = CacheEntry{Value: value, ExpiresAt: time.Now().Add(ttl)}
}

func HandleCacheStats(c *gin.Context) {
	cache.mu.RLock()
	count := len(cache.entries)
	cache.mu.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"entries": count,
		"status":  "active",
	})
}
