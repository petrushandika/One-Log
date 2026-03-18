package repository

import (
	"github.com/petrushandika/one-log/internal/domain"
	"gorm.io/gorm"
)

// Definition LogRepository interface
type LogRepository interface {
	Create(log *domain.LogEntry) error
	FindAll(limit int, offset int, sourceID, level string) ([]domain.LogEntry, int64, error)
	FindByID(id uint) (*domain.LogEntry, error)
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

func (r *logRepository) FindAll(limit int, offset int, sourceID, level string) ([]domain.LogEntry, int64, error) {
	var logs []domain.LogEntry
	var total int64
	query := r.db.Model(&domain.LogEntry{})

	if sourceID != "" {
		query = query.Where("source_id = ?", sourceID)
	}
	if level != "" {
		query = query.Where("level = ?", level)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at desc").Limit(limit).Offset(offset).Find(&logs).Error
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
