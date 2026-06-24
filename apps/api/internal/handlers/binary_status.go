package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HandleBinaryStatus(c *gin.Context) {
	binPath := os.Getenv("MIMO_BIN_PATH")
	if binPath == "" {
		binPath = "auto-detect"
	}

	_, err := os.Stat(binPath)
	exists := err == nil

	c.JSON(http.StatusOK, gin.H{
		"path":    binPath,
		"exists":  exists,
		"status":  "ready",
	})
}
