package domain

import (
"time"

"gorm.io/datatypes"
"gorm.io/gorm"
)

type Source struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"            json:"id"`
	Name      string         `gorm:"type:varchar(100);not null"          json:"name"`
	Slug      string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"slug"`
	APIKey    string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"-"`
	IsActive  bool           `gorm:"default:true"                        json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime"                      json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                               json:"-"`
}

type LogEntry struct {
	ID         uint           `gorm:"primaryKey;autoIncrement"         json:"id"`
	SourceID   string         `gorm:"index;type:varchar(50);not null"  json:"source_id"`
	Category   string         `gorm:"index;type:varchar(30);not null"  json:"category"`
	Level      string         `gorm:"type:varchar(20);not null"        json:"level"`
	Message    string         `gorm:"type:text;not null"               json:"message"`
	StackTrace string         `gorm:"type:text"                        json:"stack_trace,omitempty"`
	Context    datatypes.JSON `gorm:"type:jsonb"                       json:"context,omitempty"`
	AIInsight  datatypes.JSON `gorm:"type:jsonb"                       json:"ai_insight,omitempty"`
	IPAddress  string         `gorm:"type:varchar(45)"                 json:"ip_address,omitempty"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"                   json:"created_at"`
}
