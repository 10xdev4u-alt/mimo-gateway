package handlers

import (
	"testing"
)

func TestExtractPrompt(t *testing.T) {
	messages := []ChatMessage{
		{Role: "system", Content: "You are a helpful assistant"},
		{Role: "user", Content: "Hello world"},
		{Role: "assistant", Content: "Hi there!"},
		{Role: "user", Content: "How are you?"},
	}

	result := extractPrompt(messages)
	expected := "Hello world\nHow are you?"

	if result != expected {
		t.Errorf("extractPrompt() = %q, want %q", result, expected)
	}
}

func TestExtractPromptEmpty(t *testing.T) {
	messages := []ChatMessage{
		{Role: "system", Content: "You are a helpful assistant"},
	}

	result := extractPrompt(messages)
	if result != "" {
		t.Errorf("extractPrompt() = %q, want empty string", result)
	}
}

func TestFindMimoBinary(t *testing.T) {
	bin := findMimoBinary()
	if bin == "" {
		t.Log("Warning: mimo binary not found (expected in CI)")
	}
}
