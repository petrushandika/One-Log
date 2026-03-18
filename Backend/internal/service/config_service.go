package service

import (
	"time"

	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/internal/repository"
	"github.com/petrushandika/one-log/pkg/crypto"
)

type ConfigService interface {
	SaveConfig(sourceID, environment, key, value string, isSecret bool, updatedBy uint) error
	GetConfigsBySource(sourceID string, environment string, revealSecrets bool) ([]domain.SourceConfig, error)
	GetHistory(sourceID, environment, key string, limit int, revealSecrets bool) ([]domain.SourceConfigHistory, error)
	DeleteConfig(sourceID, environment, key string) error
}

type configService struct {
	repo repository.ConfigRepository
}

func NewConfigService(repo repository.ConfigRepository) ConfigService {
	return &configService{repo: repo}
}

func (s *configService) SaveConfig(sourceID, environment, key, value string, isSecret bool, updatedBy uint) error {
	if environment == "" {
		environment = "production"
	}
	storedValue := value
	if isSecret {
		enc, err := crypto.EncryptString(value)
		if err != nil {
			return err
		}
		storedValue = enc
	}

	now := time.Now().UTC()
	config := domain.SourceConfig{
		SourceID:    sourceID,
		Environment: environment,
		Key:         key,
		Value:       storedValue,
		IsSecret:    isSecret,
		UpdatedBy:   updatedBy,
		UpdatedAt:   now,
		CreatedAt:   now,
	}
	return s.repo.Save(&config)
}

func (s *configService) GetConfigsBySource(sourceID string, environment string, revealSecrets bool) ([]domain.SourceConfig, error) {
	configs, err := s.repo.FindBySourceID(sourceID, environment)
	if err != nil {
		return nil, err
	}

	// Mask or decrypt secrets depending on request.
	for i := range configs {
		if !configs[i].IsSecret {
			continue
		}
		if !revealSecrets {
			configs[i].Value = "****"
			continue
		}
		plain, err := crypto.DecryptString(configs[i].Value)
		if err != nil {
			return nil, err
		}
		configs[i].Value = plain
	}
	return configs, nil
}

func (s *configService) GetHistory(sourceID, environment, key string, limit int, revealSecrets bool) ([]domain.SourceConfigHistory, error) {
	h, err := s.repo.FindHistory(sourceID, environment, key, limit)
	if err != nil {
		return nil, err
	}
	for i := range h {
		if !h[i].IsSecret {
			continue
		}
		if !revealSecrets {
			h[i].Value = "****"
			continue
		}
		plain, err := crypto.DecryptString(h[i].Value)
		if err != nil {
			return nil, err
		}
		h[i].Value = plain
	}
	return h, nil
}

func (s *configService) DeleteConfig(sourceID, environment, key string) error {
	return s.repo.Delete(sourceID, environment, key)
}
