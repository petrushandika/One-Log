package service

import (
	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/internal/repository"
)

type LogService interface {
	IngestLog(req domain.IngestLogRequest, sourceID string) error
	GetLogs(limit int, page int, sourceID string, level string) ([]domain.LogEntry, int64, error)
	GetLogByID(id uint) (*domain.LogEntry, error)
}

type logService struct {
	repo repository.LogRepository
}

func NewLogService(repo repository.LogRepository) LogService {
	return &logService{repo: repo}
}

func (s *logService) IngestLog(req domain.IngestLogRequest, sourceID string) error {
	// Business Logic
	logEntry := domain.LogEntry{
		SourceID:   sourceID, // Using SourceID injected from the API Key middleware
		Category:   req.Category,
		Level:      req.Level,
		Message:    req.Message,
		StackTrace: req.StackTrace,
		IPAddress:  req.IPAddress,
	}

	// Save to DB
	err := s.repo.Create(&logEntry)
	if err != nil {
		return err
	}

	// Trigger Background Process (Goroutine)
	// Features such as: Send Email Notification on error, AI Analysis via Groq, etc.
	go func(log *domain.LogEntry) {
		// MVP Placeholder:
		if log.Level == "ERROR" || log.Level == "CRITICAL" {
			// Trigger Email Notification here
		}
	}(&logEntry)

	return nil
}

func (s *logService) GetLogs(limit int, page int, sourceID string, level string) ([]domain.LogEntry, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	return s.repo.FindAll(limit, offset, sourceID, level)
}

func (s *logService) GetLogByID(id uint) (*domain.LogEntry, error) {
	return s.repo.FindByID(id)
}

