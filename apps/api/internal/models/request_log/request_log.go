package request_log

import (
	"time"

	"gorm.io/gorm"
)

type RequestLog struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Model          string         `json:"model" gorm:"not null"`
	PromptTokens   int            `json:"prompt_tokens"`
	CompletionTokens int          `json:"completion_tokens"`
	LatencyMs      int64          `json:"latency_ms"`
	Status         string         `json:"status" gorm:"not null"`
	ErrorMessage   string         `json:"error_message,omitempty"`
	IPAddress      string         `json:"ip_address"`
	UserAgent      string         `json:"user_agent"`
	CreatedAt      time.Time      `json:"created_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (r *RequestLog) BeforeCreate(tx *gorm.DB) error {
	return nil
}
