package api_key

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

type APIKey struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Key       string         `json:"-" gorm:"uniqueIndex;not null"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Active    bool           `json:"active" gorm:"default:true"`
	LastUsed  *time.Time     `json:"last_used"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func GenerateKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "mg_" + hex.EncodeToString(bytes), nil
}

func (k *APIKey) BeforeCreate(tx *gorm.DB) error {
	if k.Key == "" {
		key, err := GenerateKey()
		if err != nil {
			return err
		}
		k.Key = key
	}
	return nil
}
