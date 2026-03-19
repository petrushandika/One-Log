package repository

import (
	"time"

	"github.com/petrushandika/one-log/internal/domain"
	"gorm.io/gorm"
)

type ActivityMonitorRepository interface {
	// Activity Feed
	GetActivityFeed(sourceID string, userID string, action string, from, to time.Time, page, limit int) ([]domain.ActivityFeed, int64, error)
	CreateActivityFeed(feed *domain.ActivityFeed) error

	// Top Active Users
	GetTopActiveUsers(sourceID string, days int, limit int) ([]domain.ActiveUser, error)

	// Activity by Resource
	GetActivityByResource(sourceID string, resourceType string, days int) ([]domain.ResourceActivity, error)

	// Before/After Diff
	StoreBeforeAfter(logID uint, before, after map[string]interface{}) error
	GetBeforeAfter(logID uint) (*domain.BeforeAfterDiff, error)

	// Compliance Export
	CreateComplianceExport(export *domain.ComplianceExport) error
	GetComplianceExports(sourceID string, page, limit int) ([]domain.ComplianceExport, int64, error)
	UpdateComplianceExportStatus(id uint, status, fileURL string) error
}

type activityMonitorRepository struct {
	db *gorm.DB
}

func NewActivityMonitorRepository(db *gorm.DB) ActivityMonitorRepository {
	return &activityMonitorRepository{db: db}
}

func (r *activityMonitorRepository) GetActivityFeed(sourceID, userID, action string, from, to time.Time, page, limit int) ([]domain.ActivityFeed, int64, error) {
	var feeds []domain.ActivityFeed
	var total int64

	query := r.db.Model(&domain.ActivityFeed{})

	if sourceID != "" {
		query = query.Where("source_id = ?", sourceID)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if !from.IsZero() {
		query = query.Where("created_at >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("created_at <= ?", to)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&feeds).Error

	return feeds, total, err
}

func (r *activityMonitorRepository) CreateActivityFeed(feed *domain.ActivityFeed) error {
	return r.db.Create(feed).Error
}

func (r *activityMonitorRepository) GetTopActiveUsers(sourceID string, days int, limit int) ([]domain.ActiveUser, error) {
	var users []domain.ActiveUser
	fromDate := time.Now().AddDate(0, 0, -days)

	query := `
		SELECT 
			user_id,
			source_id,
			COUNT(*) as activity_count,
			MAX(created_at) as last_active
		FROM activity_feeds
		WHERE created_at >= ?
	`
	args := []interface{}{fromDate}

	if sourceID != "" {
		query += " AND source_id = ?"
		args = append(args, sourceID)
	}

	query += `
		GROUP BY user_id, source_id
		ORDER BY activity_count DESC
		LIMIT ?
	`
	args = append(args, limit)

	err := r.db.Raw(query, args...).Scan(&users).Error
	return users, err
}

func (r *activityMonitorRepository) GetActivityByResource(sourceID string, resourceType string, days int) ([]domain.ResourceActivity, error) {
	var activities []domain.ResourceActivity
	fromDate := time.Now().AddDate(0, 0, -days)

	query := `
		SELECT 
			resource_type,
			resource_id,
			action,
			COUNT(*) as count,
			MAX(created_at) as last_accessed
		FROM activity_feeds
		WHERE created_at >= ?
			AND resource_type IS NOT NULL
	`
	args := []interface{}{fromDate}

	if sourceID != "" {
		query += " AND source_id = ?"
		args = append(args, sourceID)
	}

	if resourceType != "" {
		query += " AND resource_type = ?"
		args = append(args, resourceType)
	}

	query += `
		GROUP BY resource_type, resource_id, action
		ORDER BY count DESC
	`

	err := r.db.Raw(query, args...).Scan(&activities).Error
	return activities, err
}

func (r *activityMonitorRepository) StoreBeforeAfter(logID uint, before, after map[string]interface{}) error {
	// Store in a separate table or update the log entry
	// For now, we'll store it in the log entry's context
	return nil
}

func (r *activityMonitorRepository) GetBeforeAfter(logID uint) (*domain.BeforeAfterDiff, error) {
	return nil, nil
}

func (r *activityMonitorRepository) CreateComplianceExport(export *domain.ComplianceExport) error {
	return r.db.Create(export).Error
}

func (r *activityMonitorRepository) GetComplianceExports(sourceID string, page, limit int) ([]domain.ComplianceExport, int64, error) {
	var exports []domain.ComplianceExport
	var total int64

	query := r.db.Model(&domain.ComplianceExport{})
	if sourceID != "" {
		query = query.Where("source_id = ?", sourceID)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&exports).Error

	return exports, total, err
}

func (r *activityMonitorRepository) UpdateComplianceExportStatus(id uint, status, fileURL string) error {
	now := time.Now()
	return r.db.Model(&domain.ComplianceExport{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       status,
		"file_url":     fileURL,
		"completed_at": now,
	}).Error
}
