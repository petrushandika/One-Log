package domain

import "time"

// ConfigAuditTrail tracks all config changes
type ConfigAuditTrail struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	SourceID    string    `json:"source_id" gorm:"index"`
	Environment string    `json:"environment" gorm:"index"`
	Key         string    `json:"key"`
	OldValue    string    `json:"old_value"`
	NewValue    string    `json:"new_value"`
	ChangedBy   uint      `json:"changed_by" gorm:"index"`
	ChangeType  string    `json:"change_type"` // CREATE, UPDATE, DELETE
	CreatedAt   time.Time `json:"created_at"`
}

// ConfigWebhook stores webhook URLs for config change notifications
type ConfigWebhook struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	SourceID  string    `json:"source_id" gorm:"index"`
	URL       string    `json:"url"`
	Secret    string    `json:"secret"` // For HMAC signature
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
