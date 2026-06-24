package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendProxyError(c *gin.Context, status int, message string, errorType string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"message": message,
			"type":    errorType,
		},
	})
}
