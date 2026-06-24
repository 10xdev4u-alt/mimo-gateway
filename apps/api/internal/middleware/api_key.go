package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ApiKeyAuth(validKeys []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		key := strings.TrimPrefix(auth, "Bearer ")
		valid := false
		for _, k := range validKeys {
			if k == key {
				valid = true
				break
			}
		}

		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid api key"})
			c.Abort()
			return
		}

		c.Next()
	}
}
