package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Model struct {
	ID       string `json:"id"`
	Object   string `json:"object"`
	Created  int64  `json:"created"`
	OwnedBy  string `json:"owned_by"`
}

type ModelList struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

func HandleModels(c *gin.Context) {
	models := ModelList{
		Object: "list",
		Data: []Model{
			{
				ID:      "mimo-auto",
				Object:  "model",
				Created: time.Now().Unix(),
				OwnedBy: "xiaomi",
			},
			{
				ID:      "mimo-auto-stream",
				Object:  "model",
				Created: time.Now().Unix(),
				OwnedBy: "xiaomi",
			},
		},
	}
	c.JSON(http.StatusOK, models)
}
