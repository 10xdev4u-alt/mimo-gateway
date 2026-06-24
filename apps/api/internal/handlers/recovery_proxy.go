package handlers

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func RecoveryProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				_ = stack
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
