package handlers

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type ProxyRequestQueue struct {
	mu      sync.Mutex
	queue   []string
	maxSize int
}

var proxyQueue = &ProxyRequestQueue{maxSize: 100}

func (q *ProxyRequestQueue) Enqueue(item string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.queue) >= q.maxSize {
		return false
	}
	q.queue = append(q.queue, item)
	return true
}

func HandleProxyQueueStats(c *gin.Context) {
	proxyQueue.mu.Lock()
	size := len(proxyQueue.queue)
	proxyQueue.mu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"queue_size": size,
		"max_size":   proxyQueue.maxSize,
		"status":     "active",
	})
}
