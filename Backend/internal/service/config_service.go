package service

import (
	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/internal/repository"
)

type ConfigService interface {
	SaveConfig(sourceID, key, value string, isSecret bool) error
	GetConfigsBySource(sourceID string) ([]domain.SourceConfig, error)
	DeleteConfig(sourceID, key string) error
}

type configService struct {
	repo repository.ConfigRepository
}

func NewConfigService(repo repository.ConfigRepository) ConfigService {
	return &configService{repo: repo}
}

func (s *configService) SaveConfig(sourceID, key, value string, isSecret bool) error {
	config := domain.SourceConfig{
		SourceID: sourceID,
		Key:      key,
		Value:    value,
		IsSecret: isSecret,
	}
	return s.repo.Save(&config)
}

func (s *configService) GetConfigsBySource(sourceID string) ([]domain.SourceConfig, error) {
	return s.repo.FindBySourceID(sourceID)
}

func (s *configService) DeleteConfig(sourceID, key string) error {
	return s.repo.Delete(sourceID, key)
}
