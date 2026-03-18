package service

import (
	"encoding/json"

	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/internal/repository"
	"github.com/petrushandika/one-log/pkg/masking"
	"gorm.io/datatypes"
)

type LogService interface {
	IngestLog(req domain.IngestLogRequest, sourceID string) error
	GetLogs(limit int, page int, sourceID string, level string) ([]domain.LogEntry, int64, error)
	GetLogByID(id uint) (*domain.LogEntry, error)
	ManualAnalyzeLog(id uint) (*domain.LogEntry, error)
}

type logService struct {
	repo      repository.LogRepository
	notifySvc NotificationService
	aiSvc     AIService
}

func NewLogService(repo repository.LogRepository, notifySvc NotificationService, aiSvc AIService) LogService {
	return &logService{
		repo:      repo,
		notifySvc: notifySvc,
		aiSvc:     aiSvc,
	}
}

func (s *logService) IngestLog(req domain.IngestLogRequest, sourceID string) error {
	// Apply Data Masking for PII into the free-form context
	maskedContext := req.Context
	if maskedContext != nil {
		maskedContext = masking.MaskSensitiveData(maskedContext)
	}

	// Safely Marshal to JSONB datatypes
	var contextRaw datatypes.JSON
	if maskedContext != nil {
		contextBytes, _ := json.Marshal(maskedContext)
		contextRaw = datatypes.JSON(contextBytes)
	}

	// Business Logic
	logEntry := domain.LogEntry{
		SourceID:   sourceID, // Using SourceID injected from the API Key middleware
		Category:   req.Category,
		Level:      req.Level,
		Message:    req.Message,
		StackTrace: req.StackTrace,
		IPAddress:  req.IPAddress,
		Context:    contextRaw,
	}

	// Save to DB
	err := s.repo.Create(&logEntry)
	if err != nil {
		return err
	}

	// Trigger Background Process (Goroutine)
	// Features such as: Send Email Notification on error, AI Analysis via Groq, etc.
	go func(log *domain.LogEntry) {
		// 1. Trigger Notification (handles internal throttling)
		s.notifySvc.NotifyError(log)

		// 2. Automatic AI Analysis for CRITICAL logs
		if log.Level == "CRITICAL" {
			s.aiSvc.AnalyzeCriticalLog(log)
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

func (s *logService) ManualAnalyzeLog(id uint) (*domain.LogEntry, error) {
	return s.aiSvc.ManualAnalyzeLog(id)
}
