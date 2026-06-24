package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func HandleUtil36(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"util": 36, "status": "ok"})
}
