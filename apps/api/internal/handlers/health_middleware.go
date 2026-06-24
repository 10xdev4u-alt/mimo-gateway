package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
			c.Abort()
			return
		}
		c.Next()
	}
}
