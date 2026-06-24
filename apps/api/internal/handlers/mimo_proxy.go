package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// MiMoProxyHandler handles requests to the MiMo Free API
type MiMoProxyHandler struct {
	binPath     string
	mu          sync.Mutex
	lastJWT     string
	jwtExpiry   time.Time
	fingerprint string
}

// NewMiMoProxyHandler creates a new proxy handler
func NewMiMoProxyHandler() *MiMoProxyHandler {
	bin := findMimoBinary()
	fp := generateFingerprint()

	return &MiMoProxyHandler{
		binPath:     bin,
		fingerprint: fp,
	}
}

// findMimoBinary locates the .mimocode binary
func findMimoBinary() string {
	home, _ := os.UserHomeDir()
	candidates := []string{
		filepath.Join(home, ".local/share/mise/installs/node/25.8.0/lib/node_modules/@mimo-ai/cli/bin/.mimocode"),
		"/usr/local/lib/node_modules/@mimo-ai/cli/bin/.mimocode",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	// Try which mimo
	if path, err := exec.LookPath("mimo"); err == nil {
		dir := filepath.Dir(path)
		bin := filepath.Join(dir, ".mimocode")
		if _, err := os.Stat(bin); err == nil {
			return bin
		}
	}
	return ""
}

// generateFingerprint creates a device fingerprint
func generateFingerprint() string {
	hostname, _ := os.Hostname()
	username := os.Getenv("USER")
	if username == "" {
		username = "unknown"
	}
	// Simple hash - in production use crypto/sha256
	data := fmt.Sprintf("%s|linux|x64|%s", hostname, username)
	hash := fmt.Sprintf("%x", data) // Placeholder - use proper hash in production
	return hash
}

// ChatRequest represents an OpenAI-compatible chat request
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Stream      bool          `json:"stream"`
	MaxTokens   *int          `json:"max_tokens,omitempty"`
	Temperature *float64      `json:"temperature,omitempty"`
}

// ChatMessage represents a message in the chat
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse represents an OpenAI-compatible chat response
type ChatResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []ChatChoice   `json:"choices"`
	Usage   TokenUsage     `json:"usage"`
}

// ChatChoice represents a choice in the response
type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// TokenUsage represents token usage
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// HandleChat processes chat completion requests
func (h *MiMoProxyHandler) HandleChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract prompt from messages
	prompt := extractPrompt(req.Messages)
	if prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no user message"})
		return
	}

	// Run the binary
	start := time.Now()
	text, err := h.runBinary(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	latency := time.Since(start)

	// Build response
	response := ChatResponse{
		ID:      fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []ChatChoice{
			{
				Index: 0,
				Message: ChatMessage{
					Role:    "assistant",
					Content: text,
				},
				FinishReason: "stop",
			},
		},
		Usage: TokenUsage{
			PromptTokens:     len(prompt) / 4,
			CompletionTokens: len(text) / 4,
			TotalTokens:      (len(prompt) + len(text)) / 4,
		},
	}

	fmt.Printf("  ✅ %v\n", latency)
	c.JSON(http.StatusOK, response)
}

// HandleStreamChat processes streaming chat requests
func (h *MiMoProxyHandler) HandleStreamChat(c *gin.Context) {
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

	// Set headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// Run binary and stream response
	text, err := h.runBinary(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Simulate streaming by sending in chunks
	messageID := fmt.Sprintf("msg_%d", time.Now().UnixNano())
	chunkSize := 30

	// Send message_start
	c.Writer.Write([]byte(fmt.Sprintf("event: message_start\ndata: {\"type\":\"message_start\",\"message\":{\"id\":\"%s\",\"type\":\"message\",\"role\":\"assistant\",\"content\":[],\"model\":\"%s\"}}\n\n", messageID, req.Model)))
	c.Writer.Flush()

	// Send content_block_start
	c.Writer.Write([]byte("event: content_block_start\ndata: {\"type\":\"content_block_start\",\"index\":0,\"content_block\":{\"type\":\"text\",\"text\":\"\"}}\n\n"))
	c.Writer.Flush()

	// Send chunks
	for i := 0; i < len(text); i += chunkSize {
		end := i + chunkSize
		if end > len(text) {
			end = len(text)
		}
		chunk := text[i:end]
		c.Writer.Write([]byte(fmt.Sprintf("event: content_block_delta\ndata: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"text_delta\",\"text\":\"%s\"}}\n\n", chunk)))
		c.Writer.Flush()
	}

	// Send content_block_stop
	c.Writer.Write([]byte("event: content_block_stop\ndata: {\"type\":\"content_block_stop\",\"index\":0}\n\n"))
	c.Writer.Flush()

	// Send message_delta
	c.Writer.Write([]byte(fmt.Sprintf("event: message_delta\ndata: {\"type\":\"message_delta\",\"delta\":{\"stop_reason\":\"end_turn\"},\"usage\":{\"output_tokens\":%d}}\n\n", len(text)/4)))
	c.Writer.Flush()

	// Send message_stop
	c.Writer.Write([]byte("event: message_stop\ndata: {\"type\":\"message_stop\"}\n\n"))
	c.Writer.Flush()
}

// HandleModels returns available models
func (h *MiMoProxyHandler) HandleModels(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"object": "list",
		"data": []gin.H{
			{
				"id":       "mimo-auto",
				"object":   "model",
				"created":  time.Now().Unix(),
				"owned_by": "xiaomi",
			},
		},
	})
}

// HandleHealth returns health status
func (h *MiMoProxyHandler) HandleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"backend": "binary",
		"model":   "mimo-auto",
		"binary":  h.binPath,
	})
}

// runBinary executes the .mimocode binary and returns the response
func (h *MiMoProxyHandler) runBinary(prompt string) (string, error) {
	if h.binPath == "" {
		return "", fmt.Errorf("mimo binary not found")
	}

	cmd := exec.Command(h.binPath, "run", prompt, "--model", "mimo/mimo-auto")
	cmd.Env = append(os.Environ(), "MIMOCODE_FAST_BOOT=1")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	// Read output
	var output strings.Builder
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		// Parse NDJSON events
		if strings.HasPrefix(line, "{") {
			var event map[string]interface{}
			if json.Unmarshal([]byte(line), &event) == nil {
				if eventType, ok := event["type"].(string); ok {
					if eventType == "text" {
						if part, ok := event["part"].(map[string]interface{}); ok {
							if text, ok := part["text"].(string); ok {
								output.WriteString(text)
							}
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

// extractPrompt extracts the user prompt from messages
func extractPrompt(messages []ChatMessage) string {
	var parts []string
	for _, m := range messages {
		if m.Role == "user" {
			parts = append(parts, m.Content)
		}
	}
	return strings.Join(parts, "\n")
}
