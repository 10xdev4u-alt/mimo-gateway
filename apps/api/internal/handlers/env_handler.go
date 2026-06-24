package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleEnvCheck(c *gin.Context) {
	required := []string{"JWT_SECRET", "DATABASE_URL"}
	optional := []string{"REDIS_URL", "MIMO_BIN_PATH", "API_URL"}

	status := gin.H{}
	allSet := true

	for _, key := range required {
		val := getEnv(key, "")
		status[key] = gin.H{
			"set":   val != "",
			"required": true,
		}
		if val == "" {
			allSet = false
		}
	}

	for _, key := range optional {
		val := getEnv(key, "")
		status[key] = gin.H{
			"set":   val != "",
			"required": false,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"all_required_set": allSet,
		"variables":        status,
	})
}
