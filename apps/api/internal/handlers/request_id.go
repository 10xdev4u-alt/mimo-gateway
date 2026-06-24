package handlers

import (
	"fmt"
	"time"
)

func GenerateRequestID() string {
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}
