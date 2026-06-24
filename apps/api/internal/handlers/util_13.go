package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func HandleUtil13(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"util": 13, "status": "ok"})
}
