package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const Version = "1.0.0"

func HandleVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":  Version,
		"service":  "mimo-gateway",
		"compiler": "go1.24",
	})
}
