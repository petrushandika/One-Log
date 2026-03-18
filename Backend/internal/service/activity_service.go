package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/petrushandika/one-log/internal/repository"
)

type ActivityService interface {
	List(limitStr, pageStr, sourceID string, categories []string, eventType, authMethod, subjectUserID, from, to string, ownerUserID uint) (interface{}, map[string]interface{}, error)
	Summary(period, sourceID string, ownerUserID uint) (map[string]interface{}, error)
	ByUser(userID, period string, categories []string, ownerUserID uint) (map[string]interface{}, error)
	Suspicious(limitStr, pageStr, period, sourceID string, ownerUserID uint) (interface{}, map[string]interface{}, error)
}

type activityService struct {
	repo repository.LogRepository
}

func NewActivityService(repo repository.LogRepository) ActivityService {
	return &activityService{repo: repo}
}

func (s *activityService) List(limitStr, pageStr, sourceID string, categories []string, eventType, authMethod, subjectUserID, from, to string, ownerUserID uint) (interface{}, map[string]interface{}, error) {
	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	var fromT *time.Time
	var toT *time.Time
	if from != "" {
		t, err := time.Parse(time.RFC3339, from)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid from: %w", err)
		}
		fromT = &t
	}
	if to != "" {
		t, err := time.Parse(time.RFC3339, to)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid to: %w", err)
		}
		toT = &t
	}

	logs, total, err := s.repo.FindActivity(limit, offset, sourceID, categories, eventType, authMethod, subjectUserID, fromT, toT, ownerUserID)
	if err != nil {
		return nil, nil, err
	}
	meta := map[string]interface{}{
		"total": total,
		"page":  page,
		"limit": limit,
	}
	return logs, meta, nil
}

func (s *activityService) Summary(period, sourceID string, ownerUserID uint) (map[string]interface{}, error) {
	dur, err := parsePeriod(period)
	if err != nil {
		return nil, err
	}
	out, err := s.repo.GetActivitySummaryByPeriod(dur, sourceID, ownerUserID)
	if err != nil {
		return nil, err
	}
	out["period"] = period
	return out, nil
}

func (s *activityService) ByUser(userID, period string, categories []string, ownerUserID uint) (map[string]interface{}, error) {
	dur, err := parsePeriod(period)
	if err != nil {
		return nil, err
	}
	out, err := s.repo.GetUserActivity(userID, dur, categories, ownerUserID)
	if err != nil {
		return nil, err
	}
	out["period"] = period
	return out, nil
}

func (s *activityService) Suspicious(limitStr, pageStr, period, sourceID string, ownerUserID uint) (interface{}, map[string]interface{}, error) {
	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	dur, err := parsePeriod(period)
	if err != nil {
		return nil, nil, err
	}
	logs, total, err := s.repo.FindSuspiciousActivity(limit, offset, dur, sourceID, ownerUserID)
	if err != nil {
		return nil, nil, err
	}
	meta := map[string]interface{}{
		"total": total,
		"page":  page,
		"limit": limit,
	}
	return logs, meta, nil
}

func parsePeriod(p string) (time.Duration, error) {
	switch p {
	case "24h":
		return 24 * time.Hour, nil
	case "7d":
		return 7 * 24 * time.Hour, nil
	case "30d":
		return 30 * 24 * time.Hour, nil
	case "90d":
		return 90 * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unsupported period: %s", p)
	}
}
