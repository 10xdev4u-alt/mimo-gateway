package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
)

func RequestTracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Set("start_time", start)
		c.Next()
		latency := time.Since(start)
		c.Set("latency", latency)
	}
}
