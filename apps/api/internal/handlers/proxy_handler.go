package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ProxyRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

func HandleProxyChat(c *gin.Context) {
	var req ProxyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	prompt := extractPrompt(req.Messages)
	if prompt == "" {
		BadRequest(c, "no user message")
		return
	}

	proxy := NewProxyService()
	text, err := proxy.Chat(prompt)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	response := ChatResponse{
		ID:      generateID(),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []ChatChoice{{
			Index:        0,
			Message:      ChatMessage{Role: "assistant", Content: text},
			FinishReason: "stop",
		}},
		Usage: TokenUsage{
			PromptTokens:     len(prompt) / 4,
			CompletionTokens: len(text) / 4,
			TotalTokens:      (len(prompt) + len(text)) / 4,
		},
	}

	c.JSON(http.StatusOK, response)
}
