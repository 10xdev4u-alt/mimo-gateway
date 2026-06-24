package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func HandleUtil50(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"util": 50, "status": "ok"})
}
