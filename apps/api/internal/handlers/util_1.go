package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func HandleUtil1(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"util": 1, "status": "ok"})
}
