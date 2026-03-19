package repository

import (
	"github.com/petrushandika/one-log/internal/domain"
	"gorm.io/gorm"
)

// APMThresholdRepository defines methods for APM threshold management
type APMThresholdRepository interface {
	List(sourceID string) ([]domain.APMThreshold, error)
	GetByID(id uint) (*domain.APMThreshold, error)
	Create(threshold *domain.APMThreshold) error
	Update(id uint, threshold *domain.APMThreshold) error
	Delete(id uint) error
	GetByEndpoint(sourceID string, endpoint string) (*domain.APMThreshold, error)
}

type apmThresholdRepository struct {
	db *gorm.DB
}

// NewAPMThresholdRepository creates a new APM threshold repository
func NewAPMThresholdRepository(db *gorm.DB) APMThresholdRepository {
	return &apmThresholdRepository{db: db}
}

func (r *apmThresholdRepository) List(sourceID string) ([]domain.APMThreshold, error) {
	var thresholds []domain.APMThreshold
	query := r.db.Model(&domain.APMThreshold{})

	if sourceID != "" {
		query = query.Where("source_id = ?", sourceID)
	}

	if err := query.Find(&thresholds).Error; err != nil {
		return nil, err
	}
	return thresholds, nil
}

func (r *apmThresholdRepository) GetByID(id uint) (*domain.APMThreshold, error) {
	var threshold domain.APMThreshold
	err := r.db.First(&threshold, id).Error
	if err != nil {
		return nil, err
	}
	return &threshold, nil
}

func (r *apmThresholdRepository) Create(threshold *domain.APMThreshold) error {
	return r.db.Create(threshold).Error
}

func (r *apmThresholdRepository) Update(id uint, threshold *domain.APMThreshold) error {
	return r.db.Model(&domain.APMThreshold{}).Where("id = ?", id).Updates(threshold).Error
}

func (r *apmThresholdRepository) Delete(id uint) error {
	return r.db.Delete(&domain.APMThreshold{}, id).Error
}

func (r *apmThresholdRepository) GetByEndpoint(sourceID string, endpoint string) (*domain.APMThreshold, error) {
	var threshold domain.APMThreshold
	err := r.db.Where("source_id = ? AND endpoint = ?", sourceID, endpoint).First(&threshold).Error
	if err != nil {
		return nil, err
	}
	return &threshold, nil
}
