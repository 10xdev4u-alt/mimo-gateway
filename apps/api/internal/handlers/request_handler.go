package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RequestTracker struct {
	mu       sync.RWMutex
	requests []RequestInfo
}

type RequestInfo struct {
	ID        string        `json:"id"`
	Model     string        `json:"model"`
	Prompt    string        `json:"prompt"`
	Tokens    int           `json:"tokens"`
	Latency   time.Duration `json:"latency"`
	Status    string        `json:"status"`
	Timestamp time.Time     `json:"timestamp"`
}

var tracker = &RequestTracker{}

func (t *RequestTracker) Track(req RequestInfo) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.requests = append(t.requests, req)
	if len(t.requests) > 100 {
		t.requests = t.requests[len(t.requests)-100:]
	}
}

func HandleGetRequests(c *gin.Context) {
	tracker.mu.RLock()
	defer tracker.mu.RUnlock()
	c.JSON(http.StatusOK, gin.H{"requests": tracker.requests})
}
