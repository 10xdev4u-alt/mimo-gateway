package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Notification struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Read    bool   `json:"read"`
	Time    string `json:"time"`
}

func HandleGetNotifications(c *gin.Context) {
	notifications := []Notification{
		{ID: "1", Title: "Gateway Started", Message: "MiMo Gateway is running", Read: false, Time: time.Now().Format(time.RFC3339)},
		{ID: "2", Title: "New API Key", Message: "A new API key was created", Read: true, Time: time.Now().Format(time.RFC3339)},
	}
	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}

func HandleMarkNotificationRead(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"id": id, "status": "read"})
}
