package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ProxyCacheEntry struct {
	Value     interface{}
	ExpiresAt time.Time
}

type ProxySimpleCache struct {
	mu      sync.RWMutex
	entries map[string]ProxyCacheEntry
}

var proxyCache = &ProxySimpleCache{entries: make(map[string]ProxyCacheEntry)}

func (c *ProxySimpleCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	if !ok || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	return entry.Value, true
}

func (c *ProxySimpleCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = ProxyCacheEntry{Value: value, ExpiresAt: time.Now().Add(ttl)}
}

func HandleProxyCacheStats(c *gin.Context) {
	proxyCache.mu.RLock()
	count := len(proxyCache.entries)
	proxyCache.mu.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"entries": count,
		"status":  "active",
	})
}
