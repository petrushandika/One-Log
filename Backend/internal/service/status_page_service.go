package service

import (
	"time"

	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/internal/repository"
)

type StatusPageService interface {
	CreateStatusPage(sourceID, slug, title, description, logoURL string, isPublic bool) (*domain.StatusPageConfig, error)
	GetStatusPage(sourceID string) (*domain.StatusPageConfig, error)
	GetStatusPageBySlug(slug string) (*domain.StatusPageConfig, error)
	UpdateStatusPage(sourceID string, updates map[string]interface{}) error
	DeleteStatusPage(sourceID string) error
	ListStatusPages(page, limit int) ([]domain.StatusPageConfig, int64, error)

	GetUptimeStats(sourceID string, days int) (*domain.UptimeStats, error)
	GetAllUptimeStats(days int) ([]domain.UptimeStats, error)

	CreateEmbedWidget(sourceID string) (*domain.StatusPageEmbed, error)
	GetEmbedWidget(token string) (*domain.StatusPageEmbed, error)
	DeleteEmbedWidget(token string) error
}

type statusPageService struct {
	repo repository.StatusPageRepository
}

func NewStatusPageService(repo repository.StatusPageRepository) StatusPageService {
	return &statusPageService{repo: repo}
}

func (s *statusPageService) CreateStatusPage(sourceID, slug, title, description, logoURL string, isPublic bool) (*domain.StatusPageConfig, error) {
	// Check if slug already exists
	existing, _ := s.repo.GetStatusPageConfigBySlug(slug)
	if existing != nil {
		return nil, ErrSlugExists
	}

	config := &domain.StatusPageConfig{
		SourceID:    sourceID,
		Slug:        slug,
		Title:       title,
		Description: description,
		LogoURL:     logoURL,
		IsPublic:    isPublic,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreateStatusPageConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func (s *statusPageService) GetStatusPage(sourceID string) (*domain.StatusPageConfig, error) {
	return s.repo.GetStatusPageConfig(sourceID)
}

func (s *statusPageService) GetStatusPageBySlug(slug string) (*domain.StatusPageConfig, error) {
	return s.repo.GetStatusPageConfigBySlug(slug)
}

func (s *statusPageService) UpdateStatusPage(sourceID string, updates map[string]interface{}) error {
	config, err := s.repo.GetStatusPageConfig(sourceID)
	if err != nil {
		return err
	}

	// Apply updates
	if title, ok := updates["title"].(string); ok {
		config.Title = title
	}
	if desc, ok := updates["description"].(string); ok {
		config.Description = desc
	}
	if logoURL, ok := updates["logo_url"].(string); ok {
		config.LogoURL = logoURL
	}
	if isPublic, ok := updates["is_public"].(bool); ok {
		config.IsPublic = isPublic
	}

	config.UpdatedAt = time.Now()
	return s.repo.UpdateStatusPageConfig(config)
}

func (s *statusPageService) DeleteStatusPage(sourceID string) error {
	return s.repo.DeleteStatusPageConfig(sourceID)
}

func (s *statusPageService) ListStatusPages(page, limit int) ([]domain.StatusPageConfig, int64, error) {
	return s.repo.ListStatusPageConfigs(page, limit)
}

func (s *statusPageService) GetUptimeStats(sourceID string, days int) (*domain.UptimeStats, error) {
	if days <= 0 {
		days = 30
	}
	return s.repo.GetUptimeStats(sourceID, days)
}

func (s *statusPageService) GetAllUptimeStats(days int) ([]domain.UptimeStats, error) {
	if days <= 0 {
		days = 30
	}
	return s.repo.GetAllUptimeStats(days)
}

func (s *statusPageService) CreateEmbedWidget(sourceID string) (*domain.StatusPageEmbed, error) {
	return s.repo.CreateEmbedToken(sourceID)
}

func (s *statusPageService) GetEmbedWidget(token string) (*domain.StatusPageEmbed, error) {
	return s.repo.GetEmbedByToken(token)
}

func (s *statusPageService) DeleteEmbedWidget(token string) error {
	return s.repo.DeleteEmbedToken(token)
}

var ErrSlugExists = &StatusPageError{Message: "Slug already exists"}

type StatusPageError struct {
	Message string
}

func (e *StatusPageError) Error() string {
	return e.Message
}
