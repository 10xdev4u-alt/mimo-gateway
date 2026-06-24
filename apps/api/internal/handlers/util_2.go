package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func HandleUtil2(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"util": 2, "status": "ok"})
}
