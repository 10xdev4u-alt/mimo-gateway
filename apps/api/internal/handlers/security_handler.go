package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SecurityEvent struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Severity  string `json:"severity"`
	Message   string `json:"message"`
	IP        string `json:"ip"`
	Timestamp string `json:"timestamp"`
}

func HandleSecurityEvents(c *gin.Context) {
	events := []SecurityEvent{
		{ID: "1", Type: "login", Severity: "info", Message: "Successful login", IP: "192.168.1.1", Timestamp: time.Now().Format(time.RFC3339)},
		{ID: "2", Type: "rate_limit", Severity: "warning", Message: "Rate limit exceeded", IP: "10.0.0.1", Timestamp: time.Now().Format(time.RFC3339)},
	}
	c.JSON(http.StatusOK, gin.H{"events": events})
}

func HandleSecuritySummary(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"total_events": 156,
		"critical":     2,
		"warning":      23,
		"info":         131,
		"last_24h":     45,
	})
}
