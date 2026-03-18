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
	CreateSource(req domain.CreateSourceRequest) (*domain.Source, string, error)
	GetSources() ([]domain.Source, error)
	GetSourceByID(id string) (*domain.Source, error)
	RotateAPIKey(id string) (string, error)
}

type sourceService struct {
	repo repository.SourceRepository
}

func NewSourceService(repo repository.SourceRepository) SourceService {
	return &sourceService{repo: repo}
}

func (s *sourceService) CreateSource(req domain.CreateSourceRequest) (*domain.Source, string, error) {
	// Generate random API key (32 bytes = 64 hex characters)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return nil, "", err
	}
	rawAPIKey := hex.EncodeToString(bytes)
	hashedAPIKey := utils.HashAPIKey(rawAPIKey)

	source := domain.Source{
		Name:   req.Name,
		APIKey: hashedAPIKey, // Store only the SHA-256 hash in DB, never raw
	}

	if err := s.repo.Create(&source); err != nil {
		return nil, "", err
	}

	return &source, rawAPIKey, nil
}

func (s *sourceService) GetSources() ([]domain.Source, error) {
	return s.repo.FindAll()
}

func (s *sourceService) GetSourceByID(id string) (*domain.Source, error) {
	source, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if source == nil {
		return nil, errors.New("source not found")
	}
	return source, nil
}

func (s *sourceService) RotateAPIKey(id string) (string, error) {
	source, err := s.repo.FindByID(id)
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
