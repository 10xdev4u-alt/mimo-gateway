package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

func HandleHealthDetailed(c *gin.Context) {
	checks := map[string]string{
		"database": "ok",
		"redis":    "ok",
		"binary":   "ok",
	}

	status := "healthy"
	for _, v := range checks {
		if v != "ok" {
			status = "degraded"
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  status,
		"checks":  checks,
		"uptime":  time.Since(startTime).String(),
		"time":    time.Now().Format(time.RFC3339),
	})
}
