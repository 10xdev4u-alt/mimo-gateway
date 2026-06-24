package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BuildSuccessResponse(data interface{}) gin.H {
	return gin.H{
		"success": true,
		"data":    data,
	}
}

func BuildErrorResponse(message string) gin.H {
	return gin.H{
		"success": false,
		"error":   message,
	}
}

func SendSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, BuildSuccessResponse(data))
}

func SendErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, BuildErrorResponse(message))
}
