package service

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/pkg/email"
	"github.com/petrushandika/one-log/pkg/webhook"
)

// NotificationService handles alert throttling and distribution
type NotificationService interface {
	NotifyError(logEntry *domain.LogEntry)
}

type notificationService struct {
	emailClient email.SMTPEmailService
	webhook     *webhook.Client
	lastSent    sync.Map // Key: "{source_id}:{message}", Value: time.Time
}

func NewNotificationService() NotificationService {
	// Simple memory-based debounce
	return &notificationService{
		emailClient: *email.NewSMTPEmailService(),
		webhook:     webhook.New(),
	}
}

func (s *notificationService) NotifyError(logEntry *domain.LogEntry) {
	// Only proceed for ERROR or CRITICAL
	if logEntry.Level != "ERROR" && logEntry.Level != "CRITICAL" {
		return
	}

	// 1. Check Throttling (cooldown 5 minutes)
	// We create a hash/fingerprint of the issue
	issueKey := logEntry.SourceID + ":" + logEntry.Category + ":" + logEntry.Level

	val, ok := s.lastSent.Load(issueKey)
	if ok {
		lastTime, _ := val.(time.Time)
		// If the last email was sent within 5 minutes, suppress this alert
		if time.Since(lastTime) < 5*time.Minute {
			return
		}
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		return
	}

	// 2. We can send an email, first register the throttle time
	s.lastSent.Store(issueKey, time.Now())

	// 3. Send email asynchronously
	err := s.emailClient.SendAlertEmail(adminEmail, logEntry)
	if err != nil {
		log.Printf("Failed to send Notification Email to %s: %v\n", adminEmail, err)
		// Rollback throttle so it can try again
		s.lastSent.Delete(issueKey)
	} else {
		log.Printf("Notification Email Sent to %s for %s\n", adminEmail, issueKey)
	}

	// 4. Optional webhook integration (Phase 7)
	if url := os.Getenv("WEBHOOK_URL"); url != "" {
		_ = s.webhook.SendJSON(url, map[string]interface{}{
			"source_id":   logEntry.SourceID,
			"category":    logEntry.Category,
			"level":       logEntry.Level,
			"message":     logEntry.Message,
			"created_at":  logEntry.CreatedAt.UTC().Format(time.RFC3339),
			"fingerprint": logEntry.Fingerprint,
		})
	}
}
