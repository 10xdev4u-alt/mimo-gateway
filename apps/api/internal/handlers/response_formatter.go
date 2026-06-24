package handlers

import (
	"github.com/gin-gonic/gin"
)

func ResponseFormatter() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
