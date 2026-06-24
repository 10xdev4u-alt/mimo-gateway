package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleFormatJSON(c *gin.Context) {
	var req struct {
		Input string `json:"input"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, err.Error())
		return
	}

	var pretty map[string]interface{}
	if err := json.Unmarshal([]byte(req.Input), &pretty); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	output, _ := json.MarshalIndent(pretty, "", "  ")

	c.JSON(http.StatusOK, gin.H{
		"formatted": string(output),
	})
}
