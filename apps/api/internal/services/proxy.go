package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type ProxyService struct {
	binPath     string
	fingerprint string
}

func NewProxyService() *ProxyService {
	bin := findMimoBinary()
	fp := generateFingerprint()
	return &ProxyService{binPath: bin, fingerprint: fp}
}

func (s *ProxyService) Chat(prompt string) (string, error) {
	if s.binPath == "" {
		return "", fmt.Errorf("mimo binary not found")
	}

	start := time.Now()
	cmd := exec.Command(s.binPath, "run", prompt, "--model", "mimo/mimo-auto")
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

	latency := time.Since(start)
	fmt.Printf("  ✅ %v\n", latency)
	return output.String(), nil
}

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
