package service

import (
	"strconv"

	"github.com/petrushandika/one-log/internal/repository"
)

type ActivityAnalyticsService interface {
	GetAuthMethodBreakdown(sourceID string, days int) (map[string]int64, error)
	GetLoginTimeline(sourceID string, days int) ([]map[string]interface{}, error)
	GetFailedLoginHeatmap(sourceID string, days int) ([]map[string]interface{}, error)
	GetRecentSessions(limitStr, pageStr, sourceID string) (interface{}, map[string]interface{}, error)
}

type activityAnalyticsService struct {
	repo repository.ActivityAnalyticsRepository
}

func NewActivityAnalyticsService(repo repository.ActivityAnalyticsRepository) ActivityAnalyticsService {
	return &activityAnalyticsService{repo: repo}
}

func (s *activityAnalyticsService) GetAuthMethodBreakdown(sourceID string, days int) (map[string]int64, error) {
	if days <= 0 {
		days = 30
	}
	return s.repo.GetAuthMethodBreakdown(sourceID, days)
}

func (s *activityAnalyticsService) GetLoginTimeline(sourceID string, days int) ([]map[string]interface{}, error) {
	if days <= 0 {
		days = 30
	}
	return s.repo.GetLoginTimeline(sourceID, days)
}

func (s *activityAnalyticsService) GetFailedLoginHeatmap(sourceID string, days int) ([]map[string]interface{}, error) {
	if days <= 0 {
		days = 30
	}
	return s.repo.GetFailedLoginHeatmap(sourceID, days)
}

func (s *activityAnalyticsService) GetRecentSessions(limitStr, pageStr, sourceID string) (interface{}, map[string]interface{}, error) {
	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	sessions, total, err := s.repo.GetRecentSessions(limit, offset, sourceID)
	if err != nil {
		return nil, nil, err
	}

	meta := map[string]interface{}{
		"total": total,
		"page":  page,
		"limit": limit,
	}
	return sessions, meta, nil
}
