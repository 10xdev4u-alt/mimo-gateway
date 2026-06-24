package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendError(c *gin.Context, status int, message string, errorType string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"message": message,
			"type":    errorType,
		},
	})
}

func BadRequestError(c *gin.Context, message string) {
	SendError(c, http.StatusBadRequest, message, "invalid_request_error")
}

func UnauthorizedError(c *gin.Context, message string) {
	SendError(c, http.StatusUnauthorized, message, "authentication_error")
}

func NotFoundError(c *gin.Context, message string) {
	SendError(c, http.StatusNotFound, message, "not_found_error")
}

func RateLimitError(c *gin.Context) {
	SendError(c, http.StatusTooManyRequests, "rate limit exceeded", "rate_limit_error")
}

func InternalServerError(c *gin.Context, message string) {
	SendError(c, http.StatusInternalServerError, message, "server_error")
}
