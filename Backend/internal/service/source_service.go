package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/internal/repository"
	"github.com/petrushandika/one-log/pkg/utils"
)

type SourceService interface {
	CreateSource(req domain.CreateSourceRequest, userID uint) (*domain.Source, string, error)
	GetSources(userID uint) ([]domain.Source, error)
	GetSourceByID(id string, userID uint) (*domain.Source, error)
	RotateAPIKey(id string, userID uint) (string, error)
	UpdateSource(id string, userID uint, req domain.UpdateSourceRequest) (*domain.Source, error)
	DeleteSource(id string, userID uint) error
}

type sourceService struct {
	repo repository.SourceRepository
}

func NewSourceService(repo repository.SourceRepository) SourceService {
	return &sourceService{repo: repo}
}

func (s *sourceService) CreateSource(req domain.CreateSourceRequest, userID uint) (*domain.Source, string, error) {
	// Generate random API key (32 bytes = 64 hex characters)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return nil, "", err
	}
	rawAPIKey := hex.EncodeToString(bytes)
	hashedAPIKey := utils.HashAPIKey(rawAPIKey)

	source := domain.Source{
		UserID:    userID,
		Name:      req.Name,
		APIKey:    hashedAPIKey, // Store only the SHA-256 hash in DB, never raw
		HealthURL: req.HealthURL,
	}

	if err := s.repo.Create(&source); err != nil {
		return nil, "", err
	}

	return &source, rawAPIKey, nil
}

func (s *sourceService) GetSources(userID uint) ([]domain.Source, error) {
	return s.repo.FindAll(userID)
}

func (s *sourceService) GetSourceByID(id string, userID uint) (*domain.Source, error) {
	source, err := s.repo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}
	if source == nil {
		return nil, errors.New("source not found")
	}
	return source, nil
}

func (s *sourceService) RotateAPIKey(id string, userID uint) (string, error) {
	source, err := s.repo.FindByID(id, userID)
	if err != nil {
		return "", err
	}
	if source == nil {
		return "", errors.New("source not found")
	}

	// Generate new random API key
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	rawAPIKey := hex.EncodeToString(bytes)
	hashedAPIKey := utils.HashAPIKey(rawAPIKey)

	source.APIKey = hashedAPIKey

	if err := s.repo.Update(source); err != nil {
		return "", err
	}

	return rawAPIKey, nil
}

func (s *sourceService) UpdateSource(id string, userID uint, req domain.UpdateSourceRequest) (*domain.Source, error) {
	source, err := s.repo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}
	if source == nil {
		return nil, errors.New("source not found")
	}

	if req.Name != nil {
		source.Name = *req.Name
	}
	if req.HealthURL != nil {
		source.HealthURL = *req.HealthURL
	}
	if req.Status != nil {
		source.Status = *req.Status
	}

	if err := s.repo.Update(source); err != nil {
		return nil, err
	}
	return source, nil
}

func (s *sourceService) DeleteSource(id string, userID uint) error {
	// Check if source exists first
	source, err := s.repo.FindByID(id, userID)
	if err != nil {
		return err
	}
	if source == nil {
		return errors.New("source not found")
	}
	return s.repo.Delete(id, userID)
}
