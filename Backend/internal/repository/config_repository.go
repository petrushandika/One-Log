package repository

import (
	"github.com/petrushandika/one-log/internal/domain"
	"gorm.io/gorm"
)

type ConfigRepository interface {
	Save(config *domain.SourceConfig) error
	FindBySourceID(sourceID string) ([]domain.SourceConfig, error)
	Delete(sourceID, key string) error
}

type configRepository struct {
	db *gorm.DB
}

func NewConfigRepository(db *gorm.DB) ConfigRepository {
	return &configRepository{db: db}
}

func (r *configRepository) Save(config *domain.SourceConfig) error {
	var count int64
	r.db.Model(&domain.SourceConfig{}).Where("source_id = ? AND key = ?", config.SourceID, config.Key).Count(&count)
	if count > 0 {
		return r.db.Model(&domain.SourceConfig{}).Where("source_id = ? AND key = ?", config.SourceID, config.Key).Updates(config).Error
	}
	return r.db.Create(config).Error
}

func (r *configRepository) FindBySourceID(sourceID string) ([]domain.SourceConfig, error) {
	var configs []domain.SourceConfig
	err := r.db.Where("source_id = ?", sourceID).Find(&configs).Error
	return configs, err
}

func (r *configRepository) Delete(sourceID, key string) error {
	return r.db.Where("source_id = ? AND key = ?", sourceID, key).Delete(&domain.SourceConfig{}).Error
}
