package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

func HandleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "mimo-gateway",
		"version": "1.0.0",
		"uptime":  time.Since(startTime).String(),
	})
}
