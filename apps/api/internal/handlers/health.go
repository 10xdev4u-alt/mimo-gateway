package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	Binary    string `json:"binary"`
	Uptime    string `json:"uptime"`
}

var startTime = time.Now()

func HandleHealth(c *gin.Context) {
	binPath := os.Getenv("MIMO_BIN_PATH")
	if binPath == "" {
		binPath = "auto-detect"
	}

	response := HealthResponse{
		Status:  "ok",
		Service: "mimo-gateway",
		Version: "1.0.0",
		Binary:  binPath,
		Uptime:  time.Since(startTime).String(),
	}
	c.JSON(http.StatusOK, response)
}
