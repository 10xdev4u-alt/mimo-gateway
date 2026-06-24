package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleCountProxyTokens(c *gin.Context) {
	var req struct {
		Text string `json:"text"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	words := strings.Fields(req.Text)
	tokens := len(words) * 4 / 3

	c.JSON(http.StatusOK, gin.H{
		"text":   req.Text,
		"tokens": tokens,
		"words":  len(words),
		"chars":  len(req.Text),
	})
}
