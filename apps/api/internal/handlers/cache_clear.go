package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleCacheClear(c *gin.Context) {
	cleared := 0

	c.JSON(http.StatusOK, gin.H{
		"cleared": cleared,
		"status":  "ok",
	})
}
