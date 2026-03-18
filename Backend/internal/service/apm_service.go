package service

import (
	"github.com/petrushandika/one-log/internal/repository"
)

type APMService interface {
	EndpointStats(period, sourceID string, ownerUserID uint) ([]map[string]interface{}, error)
}

type apmService struct {
	repo repository.LogRepository
}

func NewAPMService(repo repository.LogRepository) APMService {
	return &apmService{repo: repo}
}

func (s *apmService) EndpointStats(period, sourceID string, ownerUserID uint) ([]map[string]interface{}, error) {
	dur, err := parsePeriod(period)
	if err != nil {
		return nil, err
	}
	rows, err := s.repo.GetAPMEndpointStats(dur, sourceID, ownerUserID)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		return []map[string]interface{}{}, nil
	}
	return rows, nil
}
