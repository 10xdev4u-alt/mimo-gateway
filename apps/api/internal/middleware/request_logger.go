package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		gin.DefaultWriter.Write([]byte(
			"[" + time.Now().Format("2006-01-02 15:04:05") + "] " +
				method + " " + path + " " +
				string(rune('0'+status/100)) + string(rune('0'+status%100/10)) + string(rune('0'+status%10)) +
				" " + latency.String() + "\n",
		))
	}
}
