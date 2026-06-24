package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleStreamProxy(c *gin.Context) {
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

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	messageID := fmt.Sprintf("msg_%d", time.Now().UnixNano())

	proxy := NewProxyService()
	text, err := proxy.Chat(prompt)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	sendSSE(c, "message_start", map[string]interface{}{
		"type": "message_start",
		"message": map[string]interface{}{
			"id": messageID, "type": "message", "role": "assistant",
			"content": []interface{}{}, "model": req.Model,
		},
	})

	sendSSE(c, "content_block_start", map[string]interface{}{
		"type": "content_block_start", "index": 0,
		"content_block": map[string]interface{}{"type": "text", "text": ""},
	})

	chunkSize := 30
	for i := 0; i < len(text); i += chunkSize {
		end := i + chunkSize
		if end > len(text) {
			end = len(text)
		}
		sendSSE(c, "content_block_delta", map[string]interface{}{
			"type": "content_block_delta", "index": 0,
			"delta": map[string]interface{}{"type": "text_delta", "text": text[i:end]},
		})
	}

	sendSSE(c, "content_block_stop", map[string]interface{}{
		"type": "content_block_stop", "index": 0,
	})

	sendSSE(c, "message_delta", map[string]interface{}{
		"type": "message_delta",
		"delta": map[string]interface{}{"stop_reason": "end_turn"},
		"usage": map[string]interface{}{"output_tokens": len(text) / 4},
	})

	sendSSE(c, "message_stop", map[string]interface{}{
		"type": "message_stop",
	})
}

func sendSSEProxy(c *gin.Context, event string, data interface{}) {
	jsonData, _ := json.Marshal(data)
	c.Writer.Write([]byte(fmt.Sprintf("event: %s\ndata: %s\n\n", event, string(jsonData))))
	c.Writer.Flush()
}
