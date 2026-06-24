package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HandleConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"app_name":     getEnv("APP_NAME", "mimo-gateway"),
		"app_env":      getEnv("APP_ENV", "development"),
		"port":         getEnv("APP_PORT", "4200"),
		"rate_limit":   getEnv("RATE_LIMIT", "100"),
		"enable_streaming": getEnv("ENABLE_STREAMING", "true"),
	})
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
