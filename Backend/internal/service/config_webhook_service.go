package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/petrushandika/one-log/internal/domain"
	"gorm.io/gorm"
)

// ConfigWebhookService handles config change notifications
type ConfigWebhookService struct {
	db *gorm.DB
}

func NewConfigWebhookService(db *gorm.DB) *ConfigWebhookService {
	return &ConfigWebhookService{db: db}
}

// SendConfigChangeNotification sends webhook notification for config changes
func (s *ConfigWebhookService) SendConfigChangeNotification(sourceID, environment, key, oldValue, newValue string) error {
	// Get webhooks for this source
	var webhooks []domain.ConfigWebhook
	if err := s.db.Where("source_id = ? AND is_active = ?", sourceID, true).Find(&webhooks).Error; err != nil {
		return err
	}

	// Prepare payload
	payload := map[string]interface{}{
		"event":       "config.changed",
		"source_id":   sourceID,
		"environment": environment,
		"key":         key,
		"old_value":   oldValue,
		"new_value":   newValue,
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Send to all registered webhooks
	for _, webhook := range webhooks {
		go s.sendWebhook(webhook, payloadBytes)
	}

	return nil
}

func (s *ConfigWebhookService) sendWebhook(webhook domain.ConfigWebhook, payload []byte) {
	req, err := http.NewRequest("POST", webhook.URL, bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ULAM-Event", "config.changed")

	// Add HMAC signature if secret is configured
	if webhook.Secret != "" {
		signature := s.generateSignature(payload, webhook.Secret)
		req.Header.Set("X-ULAM-Signature", signature)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

func (s *ConfigWebhookService) generateSignature(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

// RegisterWebhook registers a new webhook for config changes
func (s *ConfigWebhookService) RegisterWebhook(sourceID, url, secret string) (*domain.ConfigWebhook, error) {
	webhook := &domain.ConfigWebhook{
		SourceID:  sourceID,
		URL:       url,
		Secret:    secret,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(webhook).Error; err != nil {
		return nil, err
	}

	return webhook, nil
}

// ConfigPollingService handles client polling for config updates
type ConfigPollingService struct {
	db *gorm.DB
}

func NewConfigPollingService(db *gorm.DB) *ConfigPollingService {
	return &ConfigPollingService{db: db}
}

// GetConfigSince returns configs changed since given timestamp
func (s *ConfigPollingService) GetConfigSince(sourceID, environment string, since time.Time) ([]domain.SourceConfig, error) {
	var configs []domain.SourceConfig

	query := s.db.Where("source_id = ? AND updated_at > ?", sourceID, since)
	if environment != "" {
		query = query.Where("environment = ?", environment)
	}

	err := query.Find(&configs).Error
	return configs, err
}

// WatchConfig provides long-polling endpoint for config changes
func (s *ConfigPollingService) WatchConfig(sourceID, environment string, timeout time.Duration) ([]domain.SourceConfig, error) {
	// Simple implementation: poll every second until timeout or change detected
	startTime := time.Now()
	lastCheck := startTime

	for time.Since(startTime) < timeout {
		configs, err := s.GetConfigSince(sourceID, environment, lastCheck)
		if err != nil {
			return nil, err
		}

		if len(configs) > 0 {
			return configs, nil
		}

		time.Sleep(1 * time.Second)
		lastCheck = time.Now()
	}

	// Timeout - return empty
	return []domain.SourceConfig{}, nil
}
