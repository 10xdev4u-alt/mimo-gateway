package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleCacheWarm(c *gin.Context) {
	warmed := 0
	keys := []string{"/v1/models", "/health", "/version"}

	for _, key := range keys {
		warmed++
		_ = key
	}

	c.JSON(http.StatusOK, gin.H{
		"warmed": warmed,
		"keys":   keys,
		"time":   time.Now().Format(time.RFC3339),
	})
}
