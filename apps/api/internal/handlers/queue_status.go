package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleQueueStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"pending":   len(processor.items),
		"processed": processor.processed,
		"status":    "active",
	})
}
