package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func HandleUtil3(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"util": 3, "status": "ok"})
}
