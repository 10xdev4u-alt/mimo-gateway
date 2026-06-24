package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	proxyService *ProxyService
}

func NewChatHandler(proxy *ProxyService) *ChatHandler {
	return &ChatHandler{proxyService: proxy}
}

func (h *ChatHandler) HandleChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	prompt := extractPrompt(req.Messages)
	if prompt == "" {
		BadRequest(c, "no user message")
		return
	}

	text, err := h.proxyService.Chat(prompt)
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

func generateID() string {
	return "chatcmpl-" + time.Now().Format("20060102150405")
}

func extractPromptFromMessages(messages []ChatMessage) string {
	var parts []string
	for _, m := range messages {
		if m.Role == "user" {
			parts = append(parts, m.Content)
		}
	}
	return strings.Join(parts, "\n")
}
