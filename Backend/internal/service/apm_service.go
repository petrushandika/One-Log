package service

import (
	"time"

	"github.com/petrushandika/one-log/internal/repository"
)

type APMService interface {
	EndpointStats(period, sourceID string, ownerUserID uint) ([]map[string]interface{}, error)
	ResponseTimeTimeline(period, interval, sourceID, endpoint string, ownerUserID uint) ([]map[string]interface{}, error)
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

func (s *apmService) ResponseTimeTimeline(period, interval, sourceID, endpoint string, ownerUserID uint) ([]map[string]interface{}, error) {
	periodDur, err := parsePeriod(period)
	if err != nil {
		return nil, err
	}

	// Default interval is 1 hour
	intervalDur := time.Hour
	if interval != "" {
		intervalDur, err = parsePeriod(interval)
		if err != nil {
			return nil, err
		}
	}

	rows, err := s.repo.GetResponseTimeTimeline(periodDur, intervalDur, sourceID, endpoint, ownerUserID)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		return []map[string]interface{}{}, nil
	}
	return rows, nil
}
