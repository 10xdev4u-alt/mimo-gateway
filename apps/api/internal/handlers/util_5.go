package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func HandleUtil5(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"util": 5, "status": "ok"})
}
