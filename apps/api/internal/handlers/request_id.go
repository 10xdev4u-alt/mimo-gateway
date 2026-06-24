package handlers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func GenerateRequestID() string {
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}
