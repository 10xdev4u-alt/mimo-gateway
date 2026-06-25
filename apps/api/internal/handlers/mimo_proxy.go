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

type MiMoProxyHandler struct {
	binPath     string
	mu          sync.Mutex
	fingerprint string
}

func NewMiMoProxyHandler() *MiMoProxyHandler {
	bin := findMimoBinary()
	fp := generateFingerprint()
	return &MiMoProxyHandler{binPath: bin, fingerprint: fp}
}

func findMimoBinary() string {
	home, _ := os.UserHomeDir()
	candidates := []string{
		filepath.Join(home, ".local/share/mise/installs/node/25.8.0/lib/node_modules/@mimo-ai/cli/bin/.mimocode"),
		"/usr/local/bin/.mimocode",
		"/usr/local/lib/node_modules/@mimo-ai/cli/bin/.mimocode",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func generateFingerprint() string {
	hostname, _ := os.Hostname()
	username := os.Getenv("USER")
	if username == "" {
		username = "unknown"
	}
	return fmt.Sprintf("%s|%s", hostname, username)
}

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []ChatChoice `json:"choices"`
	Usage   TokenUsage   `json:"usage"`
}

type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func (h *MiMoProxyHandler) HandleChat(c *gin.Context) {
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

	start := time.Now()
	text, err := h.runBinary(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := ChatResponse{
		ID:      fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano()),
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

	fmt.Printf("  ✅ %v\n", time.Since(start))
	c.JSON(http.StatusOK, response)
}

func (h *MiMoProxyHandler) HandleModels(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"object": "list",
		"data": []gin.H{{
			"id": "mimo-auto", "object": "model",
			"created": time.Now().Unix(), "owned_by": "xiaomi",
		}},
	})
}

func (h *MiMoProxyHandler) HandleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok", "backend": "binary",
		"model": "mimo-auto", "binary": h.binPath,
	})
}

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

func extractPrompt(messages []ChatMessage) string {
	var parts []string
	for _, m := range messages {
		if m.Role == "user" {
			parts = append(parts, m.Content)
		}
	}
	return strings.Join(parts, "\n")
}
