package config

import (
	"os"
	"strconv"
	"time"
)

type GatewayConfig struct {
	Port            string
	BinPath         string
	RateLimit       int
	RateWindow      time.Duration
	AllowedOrigins  []string
	EnableStreaming  bool
	EnableLogging   bool
}

func LoadGatewayConfig() *GatewayConfig {
	return &GatewayConfig{
		Port:            getGatewayEnv("APP_PORT", "4200"),
		BinPath:         getGatewayEnv("MIMO_BIN_PATH", ""),
		RateLimit:       getEnvInt("RATE_LIMIT", 100),
		RateWindow:      getEnvDuration("RATE_WINDOW", time.Minute),
		AllowedOrigins:  getEnvList("CORS_ORIGINS", []string{"http://localhost:4201", "http://localhost:4202"}),
		EnableStreaming: getEnvBool("ENABLE_STREAMING", true),
		EnableLogging:   getEnvBool("ENABLE_LOGGING", true),
	}
}

func getGatewayEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return fallback
}

func getEnvList(key string, fallback []string) []string {
	if val := os.Getenv(key); val != "" {
		var list []string
		for _, s := range splitAndTrim(val) {
			list = append(list, s)
		}
		return list
	}
	return fallback
}

func splitAndTrim(s string) []string {
	var result []string
	for _, part := range splitString(s, ",") {
		trimmed := trimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func splitString(s, sep string) []string {
	if s == "" {
		return nil
	}
	var result []string
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
		}
	}
	result = append(result, s[start:])
	return result
}

func trimSpace(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n') {
		end--
	}
	return s[start:end]
}
