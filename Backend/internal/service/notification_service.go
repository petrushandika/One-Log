package service

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/pkg/email"
	"github.com/petrushandika/one-log/pkg/telegram"
	"github.com/petrushandika/one-log/pkg/webhook"
)

// NotificationService handles alert throttling and distribution
type NotificationService interface {
	NotifyError(logEntry *domain.LogEntry)
	NotifyRecovery(sourceName string, downtimeDuration string)
	SendEmail(to, subject, body string) error
}

type notificationService struct {
	emailClient    email.SMTPEmailService
	telegramClient *telegram.TelegramService
	webhook        *webhook.Client
	lastSent       sync.Map // Key: "{source_id}:{message}", Value: time.Time
}

func NewNotificationService() NotificationService {
	// Simple memory-based debounce
	return &notificationService{
		emailClient:    *email.NewSMTPEmailService(),
		telegramClient: telegram.NewTelegramService(),
		webhook:        webhook.New(),
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

	// 2. We can send notifications, first register the throttle time
	s.lastSent.Store(issueKey, time.Now())

	// 3. Send Email asynchronously
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail != "" {
		go func() {
			err := s.emailClient.SendAlertEmail(adminEmail, logEntry)
			if err != nil {
				log.Printf("Failed to send Notification Email to %s: %v\n", adminEmail, err)
			} else {
				log.Printf("Notification Email Sent to %s for %s\n", adminEmail, issueKey)
			}
		}()
	}

	// 4. Send Telegram asynchronously
	go func() {
		err := s.telegramClient.SendAlert(logEntry)
		if err != nil {
			log.Printf("Failed to send Telegram notification: %v\n", err)
		} else {
			log.Printf("Telegram notification sent for %s\n", issueKey)
		}
	}()

	// 5. Optional webhook integration (Phase 7)
	if url := os.Getenv("WEBHOOK_URL"); url != "" {
		go func() {
			_ = s.webhook.SendJSON(url, map[string]interface{}{
				"source_id":   logEntry.SourceID,
				"category":    logEntry.Category,
				"level":       logEntry.Level,
				"message":     logEntry.Message,
				"created_at":  logEntry.CreatedAt.UTC().Format(time.RFC3339),
				"fingerprint": logEntry.Fingerprint,
			})
		}()
	}
}

func (s *notificationService) NotifyRecovery(sourceName string, downtimeDuration string) {
	// 1. Send Email
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail != "" {
		go func() {
			err := s.emailClient.SendRecoveryEmail(adminEmail, sourceName, downtimeDuration)
			if err != nil {
				log.Printf("Failed to send Recovery Email: %v\n", err)
			} else {
				log.Printf("Recovery Email sent for %s\n", sourceName)
			}
		}()
	}

	// 2. Send Telegram
	go func() {
		err := s.telegramClient.SendRecoveryAlert(sourceName, downtimeDuration)
		if err != nil {
			log.Printf("Failed to send Telegram recovery alert: %v\n", err)
		} else {
			log.Printf("Telegram recovery alert sent for %s\n", sourceName)
		}
	}()
}

func (s *notificationService) SendEmail(to, subject, body string) error {
	return s.emailClient.SendHTML(to, subject, body)
}
