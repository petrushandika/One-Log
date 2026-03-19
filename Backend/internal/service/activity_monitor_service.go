package service

import (
	"time"

	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/internal/repository"
)

type ActivityMonitorService interface {
	GetActivityFeed(sourceID, userID, action string, from, to time.Time, page, limit int) ([]domain.ActivityFeed, int64, error)
	CreateActivityFeed(userID, sourceID, action, resourceType, resourceID string, context map[string]interface{}, ipAddress string) error
	GetTopActiveUsers(sourceID string, days, limit int) ([]domain.ActiveUser, error)
	GetActivityByResource(sourceID, resourceType string, days int) ([]domain.ResourceActivity, error)
	GetUserProfile(userID string, days int) (*UserProfile, error)
	StoreBeforeAfter(logID uint, before, after map[string]interface{}) error
	GetBeforeAfter(logID uint) (*domain.BeforeAfterDiff, error)
	RequestComplianceExport(sourceID, format string, from, to time.Time, createdBy string) (*domain.ComplianceExport, error)
	GetComplianceExports(sourceID string, page, limit int) ([]domain.ComplianceExport, int64, error)
}

type UserProfile struct {
	UserID        string                `json:"user_id"`
	TotalActivity int                   `json:"total_activity"`
	Sources       []string              `json:"sources"`
	RecentActions []domain.ActivityFeed `json:"recent_actions"`
}

type activityMonitorService struct {
	repo repository.ActivityMonitorRepository
}

func NewActivityMonitorService(repo repository.ActivityMonitorRepository) ActivityMonitorService {
	return &activityMonitorService{repo: repo}
}

func (s *activityMonitorService) GetActivityFeed(sourceID, userID, action string, from, to time.Time, page, limit int) ([]domain.ActivityFeed, int64, error) {
	return s.repo.GetActivityFeed(sourceID, userID, action, from, to, page, limit)
}

func (s *activityMonitorService) CreateActivityFeed(userID, sourceID, action, resourceType, resourceID string, context map[string]interface{}, ipAddress string) error {
	feed := &domain.ActivityFeed{
		UserID:       userID,
		SourceID:     sourceID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Context:      context,
		IPAddress:    ipAddress,
		CreatedAt:    time.Now(),
	}
	return s.repo.CreateActivityFeed(feed)
}

func (s *activityMonitorService) GetTopActiveUsers(sourceID string, days, limit int) ([]domain.ActiveUser, error) {
	if limit <= 0 {
		limit = 10
	}
	if days <= 0 {
		days = 30
	}
	return s.repo.GetTopActiveUsers(sourceID, days, limit)
}

func (s *activityMonitorService) GetActivityByResource(sourceID, resourceType string, days int) ([]domain.ResourceActivity, error) {
	if days <= 0 {
		days = 30
	}
	return s.repo.GetActivityByResource(sourceID, resourceType, days)
}

func (s *activityMonitorService) GetUserProfile(userID string, days int) (*UserProfile, error) {
	if days <= 0 {
		days = 30
	}
	fromDate := time.Now().AddDate(0, 0, -days)

	// Get all activity for this user
	feeds, _, err := s.repo.GetActivityFeed("", userID, "", fromDate, time.Time{}, 1, 100)
	if err != nil {
		return nil, err
	}

	// Collect unique sources
	sourceMap := make(map[string]bool)
	for _, feed := range feeds {
		sourceMap[feed.SourceID] = true
	}

	sources := make([]string, 0, len(sourceMap))
	for source := range sourceMap {
		sources = append(sources, source)
	}

	return &UserProfile{
		UserID:        userID,
		TotalActivity: len(feeds),
		Sources:       sources,
		RecentActions: feeds,
	}, nil
}

func (s *activityMonitorService) StoreBeforeAfter(logID uint, before, after map[string]interface{}) error {
	return s.repo.StoreBeforeAfter(logID, before, after)
}

func (s *activityMonitorService) GetBeforeAfter(logID uint) (*domain.BeforeAfterDiff, error) {
	return s.repo.GetBeforeAfter(logID)
}

func (s *activityMonitorService) RequestComplianceExport(sourceID, format string, from, to time.Time, createdBy string) (*domain.ComplianceExport, error) {
	export := &domain.ComplianceExport{
		SourceID:  sourceID,
		Format:    format,
		DateFrom:  from,
		DateTo:    to,
		Status:    "pending",
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateComplianceExport(export); err != nil {
		return nil, err
	}

	// TODO: Trigger async processing
	go s.processComplianceExport(export)

	return export, nil
}

func (s *activityMonitorService) GetComplianceExports(sourceID string, page, limit int) ([]domain.ComplianceExport, int64, error) {
	return s.repo.GetComplianceExports(sourceID, page, limit)
}

func (s *activityMonitorService) processComplianceExport(export *domain.ComplianceExport) {
	// TODO: Implement actual PDF/CSV generation
	// This is a placeholder for the async processing
}
