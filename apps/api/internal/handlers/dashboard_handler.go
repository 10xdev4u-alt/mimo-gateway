package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DashboardStats struct {
	TotalRequests   int64   `json:"total_requests"`
	AvgLatency      float64 `json:"avg_latency"`
	ActiveModels    int     `json:"active_models"`
	Uptime          string  `json:"uptime"`
	SuccessRate     float64 `json:"success_rate"`
	TotalTokens     int64   `json:"total_tokens"`
}

func HandleDashboardStats(c *gin.Context) {
	stats := DashboardStats{
		TotalRequests: 1247,
		AvgLatency:    14.2,
		ActiveModels:  1,
		Uptime:        time.Since(time.Now().Add(-2 * time.Hour)).String(),
		SuccessRate:   94.6,
		TotalTokens:   156420,
	}
	c.JSON(http.StatusOK, stats)
}

func HandleDashboardActivity(c *gin.Context) {
	activity := []gin.H{
		{"time": "14:32", "event": "request", "model": "mimo-auto", "status": "success"},
		{"time": "14:31", "event": "request", "model": "mimo-auto", "status": "success"},
		{"time": "14:30", "event": "request", "model": "mimo-auto", "status": "error"},
	}
	c.JSON(http.StatusOK, gin.H{"activity": activity})
}
