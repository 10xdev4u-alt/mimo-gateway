package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleTransformText(c *gin.Context) {
	var req struct {
		Text      string `json:"text"`
		Operation string `json:"operation"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	var result string
	switch req.Operation {
	case "uppercase":
		result = strings.ToUpper(req.Text)
	case "lowercase":
		result = strings.ToLower(req.Text)
	case "trim":
		result = strings.TrimSpace(req.Text)
	case "reverse":
		runes := []rune(req.Text)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		result = string(runes)
	default:
		result = req.Text
	}

	c.JSON(http.StatusOK, gin.H{
		"original":  req.Text,
		"result":    result,
		"operation": req.Operation,
	})
}
