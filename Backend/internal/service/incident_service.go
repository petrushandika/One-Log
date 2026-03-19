package service

import (
	"strconv"

	"github.com/petrushandika/one-log/internal/repository"
)

type IncidentService interface {
	List(limitStr, pageStr, sourceID, status string) (interface{}, map[string]interface{}, error)
	GetTimeline(sourceID string, days int) ([]map[string]interface{}, error)
}

type incidentService struct {
	repo repository.IncidentRepository
}

func NewIncidentService(repo repository.IncidentRepository) IncidentService {
	return &incidentService{repo: repo}
}

func (s *incidentService) List(limitStr, pageStr, sourceID, status string) (interface{}, map[string]interface{}, error) {
	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	incidents, total, err := s.repo.List(limit, offset, sourceID, status)
	if err != nil {
		return nil, nil, err
	}

	meta := map[string]interface{}{
		"total": total,
		"page":  page,
		"limit": limit,
	}
	return incidents, meta, nil
}

func (s *incidentService) GetTimeline(sourceID string, days int) ([]map[string]interface{}, error) {
	if days <= 0 {
		days = 30
	}
	return s.repo.GetTimeline(sourceID, days)
}
