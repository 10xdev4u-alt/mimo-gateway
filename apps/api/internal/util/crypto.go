package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateAPIKey() (string, error) {
	return GenerateRandomString(32)
}

func GenerateJWTSecret() (string, error) {
	return GenerateRandomString(64)
}

func HashString(s string) string {
	return fmt.Sprintf("%x", s)
}
