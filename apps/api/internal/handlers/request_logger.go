package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
)

type RequestLogger struct {
	logger *Logger
}

type Logger struct{}

func (l *Logger) Info(msg string, args ...interface{}) {}
func (l *Logger) Error(msg string, args ...interface{}) {}

func NewRequestLogger() *RequestLogger {
	return &RequestLogger{logger: &Logger{}}
}

func (rl *RequestLogger) Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		rl.logger.Info("%s %s %d %s", method, path, status, latency.String())
	}
}
