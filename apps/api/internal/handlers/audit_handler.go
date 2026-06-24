package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuditLog struct {
	ID        string `json:"id"`
	Action    string `json:"action"`
	User      string `json:"user"`
	Resource  string `json:"resource"`
	IP        string `json:"ip"`
	Timestamp string `json:"timestamp"`
}

func HandleAuditLogs(c *gin.Context) {
	logs := []AuditLog{
		{ID: "1", Action: "create", User: "admin", Resource: "api_key", IP: "192.168.1.1", Timestamp: time.Now().Format(time.RFC3339)},
		{ID: "2", Action: "delete", User: "admin", Resource: "api_key", IP: "192.168.1.1", Timestamp: time.Now().Format(time.RFC3339)},
	}
	c.JSON(http.StatusOK, gin.H{"logs": logs})
}
