package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func HandleUtil4(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"util": 4, "status": "ok"})
}
