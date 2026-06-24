package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequestValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > 10*1024*1024 {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "request too large",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
