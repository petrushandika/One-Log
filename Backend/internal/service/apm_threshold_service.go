package service

import (
	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/internal/repository"
)

type APMThresholdService interface {
	List(sourceID string) ([]domain.APMThreshold, error)
	GetByID(id uint) (*domain.APMThreshold, error)
	Create(sourceID string, endpoint string, p95Limit int, p99Limit int, emailNotify bool) (*domain.APMThreshold, error)
	Update(id uint, p95Limit int, p99Limit int, emailNotify bool) error
	Delete(id uint) error
	GetSlowQueries(sourceID string, thresholdMs int) ([]map[string]interface{}, error)
}

type apmThresholdService struct {
	repo    repository.APMThresholdRepository
	logRepo repository.LogRepository
}

func NewAPMThresholdService(repo repository.APMThresholdRepository, logRepo repository.LogRepository) APMThresholdService {
	return &apmThresholdService{repo: repo, logRepo: logRepo}
}

func (s *apmThresholdService) List(sourceID string) ([]domain.APMThreshold, error) {
	return s.repo.List(sourceID)
}

func (s *apmThresholdService) GetByID(id uint) (*domain.APMThreshold, error) {
	return s.repo.GetByID(id)
}

func (s *apmThresholdService) Create(sourceID string, endpoint string, p95Limit int, p99Limit int, emailNotify bool) (*domain.APMThreshold, error) {
	threshold := &domain.APMThreshold{
		SourceID:    sourceID,
		Endpoint:    endpoint,
		P95Limit:    p95Limit,
		P99Limit:    p99Limit,
		EmailNotify: emailNotify,
	}
	if err := s.repo.Create(threshold); err != nil {
		return nil, err
	}
	return threshold, nil
}

func (s *apmThresholdService) Update(id uint, p95Limit int, p99Limit int, emailNotify bool) error {
	return s.repo.Update(id, &domain.APMThreshold{
		P95Limit:    p95Limit,
		P99Limit:    p99Limit,
		EmailNotify: emailNotify,
	})
}

func (s *apmThresholdService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *apmThresholdService) GetSlowQueries(sourceID string, thresholdMs int) ([]map[string]interface{}, error) {
	if thresholdMs <= 0 {
		thresholdMs = 2000
	}
	return s.logRepo.GetSlowQueries(sourceID, thresholdMs)
}
