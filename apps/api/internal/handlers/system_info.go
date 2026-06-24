package handlers

import (
	"net/http"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
)

func HandleSystemInfo(c *gin.Context) {
	hostname, _ := os.Hostname()
	c.JSON(http.StatusOK, gin.H{
		"hostname":  hostname,
		"go_version": runtime.Version(),
		"os":        runtime.GOOS,
		"arch":      runtime.GOARCH,
		"cpus":      runtime.NumCPU(),
	})
}
