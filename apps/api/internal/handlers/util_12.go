package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func HandleUtil12(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"util": 12, "status": "ok"})
}
