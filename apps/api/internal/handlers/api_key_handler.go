package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type APIKey struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Key       string `json:"key"`
	Created   string `json:"created"`
	LastUsed  string `json:"last_used"`
	Active    bool   `json:"active"`
}

func HandleListAPIKeys(c *gin.Context) {
	keys := []APIKey{
		{ID: "1", Name: "Development", Key: "mg_abc123...", Created: "2 days ago", LastUsed: "5 min ago", Active: true},
		{ID: "2", Name: "Production", Key: "mg_def456...", Created: "1 week ago", LastUsed: "1 hour ago", Active: true},
	}
	c.JSON(http.StatusOK, gin.H{"keys": keys})
}

func HandleCreateAPIKey(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      "3",
		"name":    req.Name,
		"key":     "mg_" + time.Now().Format("20060102150405"),
		"status":  "created",
	})
}

func HandleDeleteAPIKey(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"id": id, "status": "deleted"})
}
