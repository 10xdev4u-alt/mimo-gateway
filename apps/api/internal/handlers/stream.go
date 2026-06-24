package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type StreamEvent struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func HandleStreamChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prompt := extractPrompt(req.Messages)
	if prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no user message"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	messageID := fmt.Sprintf("msg_%d", time.Now().UnixNano())
	chunkSize := 30

	text, err := runBinary(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send message_start
	sendSSE(c, "message_start", map[string]interface{}{
		"type": "message_start",
		"message": map[string]interface{}{
			"id": messageID, "type": "message", "role": "assistant",
			"content": []interface{}{}, "model": req.Model,
		},
	})

	// Send content_block_start
	sendSSE(c, "content_block_start", map[string]interface{}{
		"type": "content_block_start", "index": 0,
		"content_block": map[string]interface{}{"type": "text", "text": ""},
	})

	// Send chunks
	for i := 0; i < len(text); i += chunkSize {
		end := i + chunkSize
		if end > len(text) {
			end = len(text)
		}
		chunk := text[i:end]
		sendSSE(c, "content_block_delta", map[string]interface{}{
			"type": "content_block_delta", "index": 0,
			"delta": map[string]interface{}{"type": "text_delta", "text": chunk},
		})
	}

	// Send content_block_stop
	sendSSE(c, "content_block_stop", map[string]interface{}{
		"type": "content_block_stop", "index": 0,
	})

	// Send message_delta
	sendSSE(c, "message_delta", map[string]interface{}{
		"type": "message_delta",
		"delta": map[string]interface{}{"stop_reason": "end_turn"},
		"usage": map[string]interface{}{"output_tokens": len(text) / 4},
	})

	// Send message_stop
	sendSSE(c, "message_stop", map[string]interface{}{
		"type": "message_stop",
	})
}

func sendSSE(c *gin.Context, event string, data interface{}) {
	jsonData, _ := json.Marshal(data)
	c.Writer.Write([]byte(fmt.Sprintf("event: %s\ndata: %s\n\n", event, string(jsonData))))
	c.Writer.Flush()
}

func runBinary(prompt string) (string, error) {
	bin := findMimoBinary()
	if bin == "" {
		return "", fmt.Errorf("mimo binary not found")
	}

	cmd := exec.Command(bin, "run", prompt, "--model", "mimo/mimo-auto")
	cmd.Env = append(os.Environ(), "MIMOCODE_FAST_BOOT=1")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	var output strings.Builder
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "{") {
			var event map[string]interface{}
			if json.Unmarshal([]byte(line), &event) == nil {
				if t, ok := event["type"].(string); ok && t == "text" {
					if p, ok := event["part"].(map[string]interface{}); ok {
						if txt, ok := p["text"].(string); ok {
							output.WriteString(txt)
						}
					}
				}
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		return output.String(), err
	}

	return output.String(), nil
}
