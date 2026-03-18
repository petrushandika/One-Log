package domain

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Source presentation application or service send log
type Source struct {
	ID        string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	APIKey    string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// LogEntry presentation log from source
type LogEntry struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	SourceID   string         `gorm:"type:uuid;not null;index" json:"source_id"`
	Category   string         `gorm:"type:varchar(50);not null;index" json:"category"`
	Level      string         `gorm:"type:varchar(20);not null;index" json:"level"`
	Message    string         `gorm:"type:text;not null" json:"message"`
	Context    datatypes.JSON `gorm:"type:jsonb" json:"context"`
	StackTrace string         `gorm:"type:text" json:"stack_trace"`
	IPAddress  string         `gorm:"type:varchar(45)" json:"ip_address"`
	AIInsight  datatypes.JSON `gorm:"type:jsonb" json:"ai_insight"`
	CreatedAt  time.Time      `gorm:"index" json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
