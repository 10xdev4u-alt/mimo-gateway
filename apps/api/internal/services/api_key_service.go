package services

import (
	"time"

	"gorm.io/gorm"

	"mimo-gateway/apps/api/internal/models/api_key"
)

type APIKeyService struct {
	db *gorm.DB
}

func NewAPIKeyService(db *gorm.DB) *APIKeyService {
	return &APIKeyService{db: db}
}

func (s *APIKeyService) Create(userID uint, name string) (*api_key.APIKey, error) {
	key := &api_key.APIKey{
		Name:   name,
		UserID: userID,
		Active: true,
	}
	if err := s.db.Create(key).Error; err != nil {
		return nil, err
	}
	return key, nil
}

func (s *APIKeyService) Validate(keyStr string) (*api_key.APIKey, error) {
	var key api_key.APIKey
	if err := s.db.Where("key = ? AND active = ?", keyStr, true).First(&key).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	s.db.Model(&key).Update("last_used", &now)

	return &key, nil
}

func (s *APIKeyService) List(userID uint) ([]api_key.APIKey, error) {
	var keys []api_key.APIKey
	if err := s.db.Where("user_id = ?", userID).Find(&keys).Error; err != nil {
		return nil, err
	}
	return keys, nil
}

func (s *APIKeyService) Revoke(id uint, userID uint) error {
	return s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&api_key.APIKey{}).Error
}
