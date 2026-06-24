package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func HandleUtil21(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"util": 21, "status": "ok"})
}
