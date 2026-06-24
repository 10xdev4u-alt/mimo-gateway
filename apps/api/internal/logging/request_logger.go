package logging

import (
	"time"
	"github.com/gin-gonic/gin"
)

func RequestLogger(logger *Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		logger.Info("%s %s %d %s", method, path, status, latency.String())
	}
}
