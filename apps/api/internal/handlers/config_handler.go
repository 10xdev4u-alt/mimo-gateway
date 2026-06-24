package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HandleGetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"app_name":     getConfigEnv("APP_NAME", "mimo-gateway"),
		"app_env":      getConfigEnv("APP_ENV", "development"),
		"port":         getConfigEnv("APP_PORT", "4200"),
		"rate_limit":   getConfigEnv("RATE_LIMIT", "100"),
		"enable_streaming": getConfigEnv("ENABLE_STREAMING", "true"),
	})
}

func HandleUpdateConfig(c *gin.Context) {
	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"key":     req.Key,
		"value":   req.Value,
		"status":  "updated",
	})
}

func getConfigEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
