package handlers

import (
	"net/http"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
)

func HandleSystemHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "healthy",
		"go_version":  runtime.Version(),
		"os":          runtime.GOOS,
		"arch":        runtime.GOARCH,
		"cpus":        runtime.NumCPU(),
		"hostname":    getHostname(),
		"environment": os.Getenv("APP_ENV"),
	})
}

func getHostname() string {
	h, _ := os.Hostname()
	return h
}
