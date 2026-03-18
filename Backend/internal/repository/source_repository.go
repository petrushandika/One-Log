package repository

import (
	"errors"

	"github.com/petrushandika/one-log/internal/domain"
	"gorm.io/gorm"
)

type SourceRepository interface {
	Create(source *domain.Source) error
	FindByAPIKey(apiKey string) (*domain.Source, error)
	FindAll() ([]domain.Source, error)
	FindByID(id string) (*domain.Source, error)
	Update(source *domain.Source) error
}

type sourceRepository struct {
	db *gorm.DB
}

func NewSourceRepository(db *gorm.DB) SourceRepository {
	return &sourceRepository{db: db}
}

func (r *sourceRepository) Create(source *domain.Source) error {
	return r.db.Create(source).Error
}

func (r *sourceRepository) FindByAPIKey(apiKey string) (*domain.Source, error) {
	var source domain.Source
	err := r.db.Where("api_key = ?", apiKey).First(&source).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found but not a DB error
		}
		return nil, err
	}
	return &source, nil
}

func (r *sourceRepository) FindAll() ([]domain.Source, error) {
	var sources []domain.Source
	err := r.db.Order("created_at desc").Find(&sources).Error
	return sources, err
}

func (r *sourceRepository) FindByID(id string) (*domain.Source, error) {
	var source domain.Source
	err := r.db.Where("id = ?", id).First(&source).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &source, nil
}

func (r *sourceRepository) Update(source *domain.Source) error {
	return r.db.Save(source).Error
}
