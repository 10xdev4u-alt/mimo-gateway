package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthChecker struct {
	checks map[string]func() error
}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		checks: make(map[string]func() error),
	}
}

func (hc *HealthChecker) Register(name string, check func() error) {
	hc.checks[name] = check
}

func (hc *HealthChecker) CheckAll() map[string]string {
	results := make(map[string]string)
	for name, check := range hc.checks {
		if err := check(); err != nil {
			results[name] = "error: " + err.Error()
		} else {
			results[name] = "ok"
		}
	}
	return results
}

func HandleHealthDetailed(c *gin.Context) {
	checker := NewHealthChecker()
	results := checker.CheckAll()

	status := http.StatusOK
	for _, v := range results {
		if v != "ok" {
			status = http.StatusServiceUnavailable
			break
		}
	}

	c.JSON(status, gin.H{
		"status": "ok",
		"checks": results,
		"time":   time.Now().Format(time.RFC3339),
	})
}
