package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleValidateRequest(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid":  false,
			"error":  err.Error(),
		})
		return
	}

	errors := []string{}

	if req.Model == "" {
		errors = append(errors, "model is required")
	}
	if len(req.Messages) == 0 {
		errors = append(errors, "messages is required")
	}

	if len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid":  false,
			"errors": errors,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
	})
}
