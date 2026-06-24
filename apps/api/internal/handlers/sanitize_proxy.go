package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleSanitizeProxyInput(c *gin.Context) {
	var req struct {
		Input string `json:"input"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	sanitized := strings.TrimSpace(req.Input)
	sanitized = strings.ReplaceAll(sanitized, "<script>", "")
	sanitized = strings.ReplaceAll(sanitized, "</script>", "")

	c.JSON(http.StatusOK, gin.H{
		"original":  req.Input,
		"sanitized": sanitized,
	})
}
