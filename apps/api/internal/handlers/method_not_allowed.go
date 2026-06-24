package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MethodNotAllowedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error":  "method not allowed",
			"method": c.Request.Method,
		})
	}
}
