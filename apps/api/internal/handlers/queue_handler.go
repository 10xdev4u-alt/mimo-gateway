package handlers

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type RequestQueue struct {
	mu      sync.Mutex
	queue   []string
	maxSize int
}

var requestQueue = &RequestQueue{maxSize: 100}

func (q *RequestQueue) Enqueue(item string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.queue) >= q.maxSize {
		return false
	}
	q.queue = append(q.queue, item)
	return true
}

func (q *RequestQueue) Dequeue() (string, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.queue) == 0 {
		return "", false
	}
	item := q.queue[0]
	q.queue = q.queue[1:]
	return item, true
}

func HandleQueueStats(c *gin.Context) {
	requestQueue.mu.Lock()
	size := len(requestQueue.queue)
	requestQueue.mu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"queue_size": size,
		"max_size":   requestQueue.maxSize,
		"status":     "active",
	})
}
