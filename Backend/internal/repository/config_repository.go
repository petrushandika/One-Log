package repository

import (
	"github.com/petrushandika/one-log/internal/domain"
	"gorm.io/gorm"
)

type ConfigRepository interface {
	Save(config *domain.SourceConfig) error
	FindBySourceID(sourceID string, environment string) ([]domain.SourceConfig, error)
	FindHistory(sourceID, environment, key string, limit int) ([]domain.SourceConfigHistory, error)
	Delete(sourceID, environment, key string) error
}

type configRepository struct {
	db *gorm.DB
}

func NewConfigRepository(db *gorm.DB) ConfigRepository {
	return &configRepository{db: db}
}

func (r *configRepository) Save(config *domain.SourceConfig) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing domain.SourceConfig
		err := tx.Where("source_id = ? AND environment = ? AND key = ?", config.SourceID, config.Environment, config.Key).First(&existing).Error
		version := int64(1)
		if err == nil {
			// Update existing row
			config.ID = existing.ID
			version = 0
			// Compute next version from history table
			tx.Model(&domain.SourceConfigHistory{}).
				Where("source_id = ? AND environment = ? AND key = ?", config.SourceID, config.Environment, config.Key).
				Select("COALESCE(MAX(version),0)+1").
				Scan(&version)
			if err := tx.Model(&domain.SourceConfig{}).
				Where("id = ?", existing.ID).
				Updates(config).Error; err != nil {
				return err
			}
		} else {
			// Create new row
			tx.Model(&domain.SourceConfigHistory{}).
				Where("source_id = ? AND environment = ? AND key = ?", config.SourceID, config.Environment, config.Key).
				Select("COALESCE(MAX(version),0)+1").
				Scan(&version)
			if err := tx.Create(config).Error; err != nil {
				return err
			}
		}

		h := domain.SourceConfigHistory{
			SourceID:    config.SourceID,
			Environment: config.Environment,
			Key:         config.Key,
			Value:       config.Value,
			IsSecret:    config.IsSecret,
			Version:     version,
			UpdatedBy:   config.UpdatedBy,
			CreatedAt:   config.UpdatedAt,
		}
		return tx.Create(&h).Error
	})
}

func (r *configRepository) FindBySourceID(sourceID string, environment string) ([]domain.SourceConfig, error) {
	var configs []domain.SourceConfig
	query := r.db.Where("source_id = ?", sourceID)
	if environment != "" {
		query = query.Where("environment = ?", environment)
	}
	err := query.Order("key asc").Find(&configs).Error
	return configs, err
}

func (r *configRepository) FindHistory(sourceID, environment, key string, limit int) ([]domain.SourceConfigHistory, error) {
	var out []domain.SourceConfigHistory
	if limit <= 0 {
		limit = 50
	}
	err := r.db.Where("source_id = ? AND environment = ? AND key = ?", sourceID, environment, key).
		Order("version desc").
		Limit(limit).
		Find(&out).Error
	return out, err
}

func (r *configRepository) Delete(sourceID, environment, key string) error {
	return r.db.Where("source_id = ? AND environment = ? AND key = ?", sourceID, environment, key).Delete(&domain.SourceConfig{}).Error
}
