package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequestSizeMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "request too large",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
