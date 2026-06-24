package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Model struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	MaxTokens   int     `json:"max_tokens"`
	Cost        float64 `json:"cost"`
}

func HandleListModels(c *gin.Context) {
	models := []Model{
		{ID: "mimo-auto", Name: "MiMo Auto", Description: "Default model", MaxTokens: 128000, Cost: 0},
	}
	c.JSON(http.StatusOK, gin.H{"models": models})
}

func HandleGetModel(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, Model{
		ID:          id,
		Name:        "MiMo Auto",
		Description: "Default model",
		MaxTokens:   128000,
		Cost:        0,
		CreatedAt:   time.Now().Format(time.RFC3339),
	})
}
