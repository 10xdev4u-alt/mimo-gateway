package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleRequestHistory(c *gin.Context) {
	history := []gin.H{
		{"id": "1", "model": "mimo-auto", "tokens": 156, "latency": "12.4s", "time": "2m ago", "status": "success"},
		{"id": "2", "model": "mimo-auto", "tokens": 89, "latency": "14.8s", "time": "5m ago", "status": "success"},
		{"id": "3", "model": "mimo-auto", "tokens": 234, "latency": "16.2s", "time": "8m ago", "status": "success"},
	}
	c.JSON(http.StatusOK, gin.H{"data": history})
}
