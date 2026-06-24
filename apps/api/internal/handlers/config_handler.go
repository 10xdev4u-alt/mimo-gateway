package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HandleConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"app_name":     getConfigEnv("APP_NAME", "mimo-gateway"),
		"app_env":      getConfigEnv("APP_ENV", "development"),
		"port":         getConfigEnv("APP_PORT", "4200"),
		"rate_limit":   getConfigEnv("RATE_LIMIT", "100"),
		"enable_streaming": getConfigEnv("ENABLE_STREAMING", "true"),
	})
}

func getConfigEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
