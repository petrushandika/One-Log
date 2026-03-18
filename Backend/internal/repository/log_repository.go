package repository

import (
	"github.com/petrushandika/one-log/internal/domain"
	"gorm.io/gorm"
)

// Definition LogRepository interface
type LogRepository interface {
	Create(log *domain.LogEntry) error
	FindAll(limit int, offset int, sourceID, level, category string, userID uint) ([]domain.LogEntry, int64, error)
	FindByID(id uint) (*domain.LogEntry, error)
	Update(log *domain.LogEntry) error
	DeleteOlderThan(days int) error
	GetStatsOverview(userID uint) (map[string]interface{}, error)
	CountFailedAttempts(ip string, durationMinutes int) (int64, error)
}

// Struct private for implementation
type logRepository struct {
	db *gorm.DB
}

// Constructor
func NewLogRepository(db *gorm.DB) LogRepository {
	return &logRepository{db: db}
}

// Implementation method Create
func (r *logRepository) Create(log *domain.LogEntry) error {
	return r.db.Create(log).Error
}

func (r *logRepository) FindAll(limit int, offset int, sourceID, level, category string, userID uint) ([]domain.LogEntry, int64, error) {
	var logs []domain.LogEntry
	var total int64
	query := r.db.Model(&domain.LogEntry{})

	if userID > 0 {
		query = query.Joins("JOIN sources ON sources.id = log_entries.source_id").Where("sources.user_id = ?", userID)
	}

	if sourceID != "" {
		query = query.Where("log_entries.source_id = ?", sourceID)
	}
	if level != "" {
		query = query.Where("log_entries.level = ?", level)
	}
	if category != "" {
		query = query.Where("log_entries.category = ?", category)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("log_entries.created_at desc").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}

func (r *logRepository) FindByID(id uint) (*domain.LogEntry, error) {
	var log domain.LogEntry
	err := r.db.First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *logRepository) Update(log *domain.LogEntry) error {
	return r.db.Save(log).Error
}

func (r *logRepository) DeleteOlderThan(days int) error {
	// Fase 2 Immutable Logs Guard:
	// Logs categorized as AUDIT_TRAIL should NOT be deleted by general retention policies.
	return r.db.Where("created_at < NOW() - INTERVAL '1 day' * ? AND category != 'AUDIT_TRAIL'", days).Delete(&domain.LogEntry{}).Error
}

func (r *logRepository) GetStatsOverview(userID uint) (map[string]interface{}, error) {
	queryTotal := r.db.Model(&domain.LogEntry{})
	queryBreakdown := r.db.Model(&domain.LogEntry{})

	if userID > 0 {
		queryTotal = queryTotal.Joins("JOIN sources ON sources.id = log_entries.source_id").Where("sources.user_id = ?", userID)
		queryBreakdown = queryBreakdown.Joins("JOIN sources ON sources.id = log_entries.source_id").Where("sources.user_id = ?", userID)
	}

	var total int64
	if err := queryTotal.Count(&total).Error; err != nil {
		return nil, err
	}

	type result struct {
		Level string
		Count int64
	}
	var bResult []result
	if err := queryBreakdown.Select("log_entries.level, count(*) as count").Group("log_entries.level").Find(&bResult).Error; err != nil {
		return nil, err
	}

	stats := map[string]interface{}{"total": total}
	for _, row := range bResult {
		stats[row.Level] = row.Count
	}

	return stats, nil
}

func (r *logRepository) CountFailedAttempts(ip string, durationMinutes int) (int64, error) {
	var count int64
	err := r.db.Model(&domain.LogEntry{}).
		Where("ip_address = ? AND category = 'AUTH_EVENT' AND level = 'WARN' AND created_at >= NOW() - INTERVAL '1 minute' * ?", ip, durationMinutes).
		Count(&count).Error
	return count, err
}
