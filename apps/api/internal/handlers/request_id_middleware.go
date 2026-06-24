package handlers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := fmt.Sprintf("req_%d", time.Now().UnixNano())
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}
