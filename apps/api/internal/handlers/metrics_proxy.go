package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ProxyMetrics struct {
	mu            sync.RWMutex
	totalRequests int64
	totalTokens   int64
	totalLatency  time.Duration
	errorCount    int64
}

var proxyMetrics = &ProxyMetrics{}

func (m *ProxyMetrics) RecordRequest(tokens int, latency time.Duration, isError bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.totalRequests++
	m.totalTokens += int64(tokens)
	m.totalLatency += latency
	if isError {
		m.errorCount++
	}
}

func HandleGetProxyMetrics(c *gin.Context) {
	proxyMetrics.mu.RLock()
	defer proxyMetrics.mu.RUnlock()

	var avgLatency float64
	if proxyMetrics.totalRequests > 0 {
		avgLatency = float64(proxyMetrics.totalLatency.Milliseconds()) / float64(proxyMetrics.totalRequests)
	}

	var errorRate float64
	if proxyMetrics.totalRequests > 0 {
		errorRate = float64(proxyMetrics.errorCount) / float64(proxyMetrics.totalRequests) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"total_requests": proxyMetrics.totalRequests,
		"total_tokens":   proxyMetrics.totalTokens,
		"avg_latency_ms": avgLatency,
		"error_rate":     errorRate,
	})
}
