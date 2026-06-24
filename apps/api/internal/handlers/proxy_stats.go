package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ProxyStats struct {
	mu            sync.RWMutex
	requests      []RequestStat
	totalTokens   int
	totalLatency  time.Duration
}

type RequestStat struct {
	ID        string
	Model     string
	Tokens    int
	Latency   time.Duration
	Timestamp time.Time
	Status    string
}

var proxyStats = &ProxyStats{}

func (ps *ProxyStats) Record(stat RequestStat) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.requests = append(ps.requests, stat)
	ps.totalTokens += stat.Tokens
	ps.totalLatency += stat.Latency
	if len(ps.requests) > 100 {
		ps.requests = ps.requests[len(ps.requests)-100:]
	}
}

func (ps *ProxyStats) GetRecent(limit int) []RequestStat {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	if limit > len(ps.requests) {
		limit = len(ps.requests)
	}
	return ps.requests[len(ps.requests)-limit:]
}

func HandleProxyStats(c *gin.Context) {
	stats := proxyStats.GetRecent(10)
	c.JSON(http.StatusOK, gin.H{
		"recent_requests": stats,
		"total_tokens":    proxyStats.totalTokens,
	})
}
