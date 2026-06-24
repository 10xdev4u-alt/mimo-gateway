package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LogEntry struct {
	ID        string `json:"id"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Source    string `json:"source"`
	Timestamp string `json:"timestamp"`
}

func HandleGetLogs(c *gin.Context) {
	logs := []LogEntry{
		{ID: "1", Level: "info", Message: "Request processed", Source: "proxy", Timestamp: time.Now().Format(time.RFC3339)},
		{ID: "2", Level: "error", Message: "Binary timeout", Source: "proxy", Timestamp: time.Now().Format(time.RFC3339)},
	}
	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

func HandleGetLogsByLevel(c *gin.Context) {
	level := c.Param("level")
	c.JSON(http.StatusOK, gin.H{"level": level, "count": 42})
}
