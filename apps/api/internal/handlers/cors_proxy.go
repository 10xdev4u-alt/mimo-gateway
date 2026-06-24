package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleProxyCORSPreflight(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Header("Access-Control-Max-Age", "86400")

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
