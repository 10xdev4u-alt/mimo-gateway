package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Param   interface{} `json:"param"`
	Code    interface{} `json:"code"`
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: ErrorDetail{
			Message: message,
			Type:    "invalid_request_error",
		},
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Error: ErrorDetail{
			Message: message,
			Type:    "authentication_error",
		},
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Error: ErrorDetail{
			Message: message,
			Type:    "not_found_error",
		},
	})
}

func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: ErrorDetail{
			Message: message,
			Type:    "server_error",
		},
	})
}

func RateLimited(c *gin.Context) {
	c.JSON(http.StatusTooManyRequests, ErrorResponse{
		Error: ErrorDetail{
			Message: "rate limit exceeded",
			Type:    "rate_limit_error",
		},
	})
}
