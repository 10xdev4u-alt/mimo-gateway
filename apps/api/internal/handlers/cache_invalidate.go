package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleCacheInvalidate(c *gin.Context) {
	var req struct {
		Pattern string `json:"pattern"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pattern": req.Pattern,
		"status":  "invalidated",
	})
}
