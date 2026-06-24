package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type MetricsCollector struct {
	mu            sync.RWMutex
	totalRequests int64
	totalTokens   int64
	totalLatency  time.Duration
	errorCount    int64
}

var metrics = &MetricsCollector{}

func (m *MetricsCollector) RecordRequest(tokens int, latency time.Duration, isError bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.totalRequests++
	m.totalTokens += tokens
	m.totalLatency += latency
	if isError {
		m.errorCount++
	}
}

func HandleGetMetrics(c *gin.Context) {
	metrics.mu.RLock()
	defer metrics.mu.RUnlock()

	var avgLatency float64
	if metrics.totalRequests > 0 {
		avgLatency = float64(metrics.totalLatency.Milliseconds()) / float64(metrics.totalRequests)
	}

	var errorRate float64
	if metrics.totalRequests > 0 {
		errorRate = float64(metrics.errorCount) / float64(metrics.totalRequests) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"total_requests": metrics.totalRequests,
		"total_tokens":   metrics.totalTokens,
		"avg_latency_ms": avgLatency,
		"error_rate":     errorRate,
	})
}
