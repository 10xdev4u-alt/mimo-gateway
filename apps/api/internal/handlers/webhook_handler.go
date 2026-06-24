package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleWebhook(c *gin.Context) {
	provider := c.Param("provider")

	c.JSON(http.StatusOK, gin.H{
		"provider": provider,
		"status":   "received",
		"message":  "webhook processed",
	})
}
