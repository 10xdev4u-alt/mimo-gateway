package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AnalyticsData struct {
	RequestsToday    int     `json:"requests_today"`
	RequestsYesterday int    `json:"requests_yesterday"`
	TokensUsed       int64   `json:"tokens_used"`
	AvgResponseTime  float64 `json:"avg_response_time"`
	ErrorRate        float64 `json:"error_rate"`
}

func HandleAnalytics(c *gin.Context) {
	data := AnalyticsData{
		RequestsToday:     245,
		RequestsYesterday: 312,
		TokensUsed:        31250,
		AvgResponseTime:   14.2,
		ErrorRate:         2.4,
	}
	c.JSON(http.StatusOK, data)
}

func HandleAnalyticsTimeline(c *gin.Context) {
	timeline := []gin.H{
		{"hour": "00:00", "requests": 12, "tokens": 1500},
		{"hour": "06:00", "requests": 45, "tokens": 5625},
		{"hour": "12:00", "requests": 89, "tokens": 11125},
		{"hour": "18:00", "requests": 67, "tokens": 8375},
	}
	c.JSON(http.StatusOK, gin.H{
		"timeline": timeline,
		"period":   "24h",
	})
}
