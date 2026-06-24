package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFoundProxyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "endpoint not found",
			"path":  c.Request.URL.Path,
		})
	}
}
