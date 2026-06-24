package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func HandleUtil25(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"util": 25, "status": "ok"})
}
