package service

import (
	"fmt"
	"strconv"

	"github.com/petrushandika/one-log/internal/repository"
)

type IssueService interface {
	List(limitStr, pageStr, sourceID, status string, ownerUserID uint) (interface{}, map[string]interface{}, error)
	Get(fingerprint string, ownerUserID uint) (interface{}, error)
	UpdateStatus(fingerprint string, status string, ownerUserID uint) (interface{}, error)
	Logs(fingerprint string, limitStr, pageStr string, ownerUserID uint) (interface{}, map[string]interface{}, error)
}

type issueService struct {
	repo repository.LogRepository
}

func NewIssueService(repo repository.LogRepository) IssueService {
	return &issueService{repo: repo}
}

func (s *issueService) List(limitStr, pageStr, sourceID, status string, ownerUserID uint) (interface{}, map[string]interface{}, error) {
	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	issues, total, err := s.repo.ListIssues(limit, offset, sourceID, status, ownerUserID)
	if err != nil {
		return nil, nil, err
	}
	meta := map[string]interface{}{
		"total": total,
		"page":  page,
		"limit": limit,
	}
	return issues, meta, nil
}

func (s *issueService) Get(fingerprint string, ownerUserID uint) (interface{}, error) {
	issue, err := s.repo.GetIssueByFingerprint(fingerprint, ownerUserID)
	if err != nil {
		return nil, err
	}
	return issue, nil
}

func (s *issueService) UpdateStatus(fingerprint string, status string, ownerUserID uint) (interface{}, error) {
	if status != "OPEN" && status != "RESOLVED" && status != "IGNORED" {
		return nil, fmt.Errorf("invalid status: %s", status)
	}
	issue, err := s.repo.UpdateIssueStatus(fingerprint, status, ownerUserID)
	if err != nil {
		return nil, err
	}
	return issue, nil
}

func (s *issueService) Logs(fingerprint string, limitStr, pageStr string, ownerUserID uint) (interface{}, map[string]interface{}, error) {
	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	logs, total, err := s.repo.ListIssueLogs(limit, offset, fingerprint, ownerUserID)
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
