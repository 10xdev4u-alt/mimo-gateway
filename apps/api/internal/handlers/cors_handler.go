package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleCORSPreflight(c *gin.Context) {
	origin := c.GetHeader("Origin")
	allowedOrigins := []string{"http://localhost:4201", "http://localhost:4202"}

	allowed := false
	for _, o := range allowedOrigins {
		if strings.EqualFold(o, origin) {
			allowed = true
			break
		}
	}

	if allowed {
		c.Header("Access-Control-Allow-Origin", origin)
	}
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Header("Access-Control-Max-Age", "86400")

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
