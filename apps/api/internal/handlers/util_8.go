package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func HandleUtil8(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"util": 8, "status": "ok"})
}
