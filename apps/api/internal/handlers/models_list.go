package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleModelsList(c *gin.Context) {
	models := []gin.H{
		{
			"id":       "mimo-auto",
			"object":   "model",
			"created":  time.Now().Unix(),
			"owned_by": "xiaomi",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"object": "list",
		"data":   models,
	})
}
