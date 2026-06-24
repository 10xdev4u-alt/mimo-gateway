package handlers

import (
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func HandleBinaryStatus(c *gin.Context) {
	binPath := os.Getenv("MIMO_BIN_PATH")
	if binPath == "" {
		binPath = "auto-detect"
	}

	installed := false
	if _, err := os.Stat(binPath); err == nil {
		installed = true
	}

	version := ""
	if installed {
		cmd := exec.Command(binPath, "--version")
		if out, err := cmd.Output(); err == nil {
			version = string(out)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"path":      binPath,
		"installed": installed,
		"version":   version,
	})
}
