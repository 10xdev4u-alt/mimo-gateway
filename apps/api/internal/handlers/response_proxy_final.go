package handlers

import (
	"github.com/gin-gonic/gin"
)

func ResponseProxyFormatter() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
