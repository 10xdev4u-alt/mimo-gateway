package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Metrics struct {
	mu            sync.RWMutex
	totalRequests int64
	errorCount    int64
	totalLatency  time.Duration
	startTime     time.Time
}

var globalMetrics = &Metrics{startTime: time.Now()}

func (m *Metrics) RecordRequest(latency time.Duration, isError bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.totalRequests++
	m.totalLatency += latency
	if isError {
		m.errorCount++
	}
}

func (m *Metrics) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var avgLatency float64
	if m.totalRequests > 0 {
		avgLatency = float64(m.totalLatency.Milliseconds()) / float64(m.totalRequests)
	}

	return map[string]interface{}{
		"total_requests": m.totalRequests,
		"error_count":    m.errorCount,
		"avg_latency_ms": avgLatency,
		"uptime":         time.Since(m.startTime).String(),
	}
}

func HandleMetrics(c *gin.Context) {
	stats := globalMetrics.GetStats()
	c.JSON(http.StatusOK, stats)
}
