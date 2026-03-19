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
	GetSlowQueryTrend(sourceID string, days int) ([]map[string]interface{}, error)
	CalculateApdexScore(sourceID string, endpoint string, thresholdMs int) (*ApdexScore, error)
	CheckThresholdAlerts() ([]ThresholdAlert, error)
}

type ApdexScore struct {
	SourceID   string  `json:"source_id"`
	Endpoint   string  `json:"endpoint"`
	Score      float64 `json:"score"`
	Satisfied  int     `json:"satisfied"`
	Tolerating int     `json:"tolerating"`
	Frustrated int     `json:"frustrated"`
	Total      int     `json:"total"`
}

type ThresholdAlert struct {
	Threshold  domain.APMThreshold `json:"threshold"`
	CurrentP95 int                 `json:"current_p95"`
	CurrentP99 int                 `json:"current_p99"`
	Exceeded   bool                `json:"exceeded"`
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

func (s *apmThresholdService) GetSlowQueryTrend(sourceID string, days int) ([]map[string]interface{}, error) {
	if days <= 0 {
		days = 30
	}
	return s.logRepo.GetSlowQueryTrend(sourceID, days)
}

func (s *apmThresholdService) CalculateApdexScore(sourceID string, endpoint string, thresholdMs int) (*ApdexScore, error) {
	if thresholdMs <= 0 {
		thresholdMs = 1000 // Default threshold
	}

	result, err := s.logRepo.CalculateApdexScore(sourceID, endpoint, thresholdMs)
	if err != nil {
		return nil, err
	}

	return &ApdexScore{
		SourceID:   sourceID,
		Endpoint:   endpoint,
		Score:      result.Score,
		Satisfied:  result.Satisfied,
		Tolerating: result.Tolerating,
		Frustrated: result.Frustrated,
		Total:      result.Total,
	}, nil
}

func (s *apmThresholdService) CheckThresholdAlerts() ([]ThresholdAlert, error) {
	// Get all thresholds with email notification enabled
	thresholds, err := s.repo.List("")
	if err != nil {
		return nil, err
	}

	var alerts []ThresholdAlert
	for _, threshold := range thresholds {
		if !threshold.EmailNotify {
			continue
		}

		// Get current latency stats
		stats, err := s.logRepo.GetEndpointLatencyStats(threshold.SourceID, threshold.Endpoint)
		if err != nil {
			continue
		}

		p95 := int(stats["p95_ms"].(float64))
		p99 := int(stats["p99_ms"].(float64))

		exceeded := p95 > threshold.P95Limit || p99 > threshold.P99Limit

		if exceeded {
			alerts = append(alerts, ThresholdAlert{
				Threshold:  threshold,
				CurrentP95: p95,
				CurrentP99: p99,
				Exceeded:   true,
			})
		}
	}

	return alerts, nil
}
