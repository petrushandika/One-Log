package domain

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Source struct {
	ID        string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"` // Owner of this source
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	APIKey    string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"-"`
	HealthURL string         `gorm:"type:varchar(255)" json:"health_url"`
	Status    string         `gorm:"type:varchar(20);default:'ONLINE'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type LogEntry struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	SourceID    string         `gorm:"type:uuid;not null;index" json:"source_id"`
	Category    string         `gorm:"type:varchar(50);not null;index" json:"category"`
	Level       string         `gorm:"type:varchar(20);not null;index" json:"level"`
	Message     string         `gorm:"type:text;not null" json:"message"`
	Context     datatypes.JSON `gorm:"type:jsonb" json:"context"`
	StackTrace  string         `gorm:"type:text" json:"stack_trace"`
	IPAddress   string         `gorm:"type:varchar(45)" json:"ip_address"`
	AIInsight   datatypes.JSON `gorm:"type:jsonb" json:"ai_insight"`
	Fingerprint string         `gorm:"type:varchar(64);index" json:"fingerprint"`
	CreatedAt   time.Time      `gorm:"index" json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Issue is an aggregated error group keyed by LogEntry.Fingerprint.
// Phase 5: Issue Tracker
type Issue struct {
	Fingerprint     string    `gorm:"type:varchar(64);primaryKey" json:"fingerprint"`
	SourceID        string    `gorm:"type:uuid;not null;index" json:"source_id"`
	Status          string    `gorm:"type:varchar(20);default:'OPEN';index" json:"status"`
	Category        string    `gorm:"type:varchar(50);index" json:"category"`
	Level           string    `gorm:"type:varchar(20);index" json:"level"`
	MessageSample   string    `gorm:"type:text" json:"message_sample"`
	OccurrenceCount int64     `gorm:"not null;default:1" json:"occurrence_count"`
	FirstSeenAt     time.Time `gorm:"index" json:"first_seen_at"`
	LastSeenAt      time.Time `gorm:"index" json:"last_seen_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type SourceConfig struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SourceID    string    `gorm:"type:uuid;not null;index" json:"source_id"`
	Environment string    `gorm:"type:varchar(30);not null;default:'production';index" json:"environment"`
	Key         string    `gorm:"type:varchar(100);not null;index" json:"key"`
	Value       string    `gorm:"type:text" json:"value"`
	IsSecret    bool      `gorm:"default:false" json:"is_secret"`
	UpdatedBy   uint      `gorm:"index" json:"updated_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// SourceConfigHistory stores immutable config revisions for audit/rollback.
type SourceConfigHistory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SourceID    string    `gorm:"type:uuid;not null;index" json:"source_id"`
	Environment string    `gorm:"type:varchar(30);not null;index" json:"environment"`
	Key         string    `gorm:"type:varchar(100);not null;index" json:"key"`
	Value       string    `gorm:"type:text" json:"value"`
	IsSecret    bool      `gorm:"default:false" json:"is_secret"`
	Version     int64     `gorm:"not null;default:1;index" json:"version"`
	UpdatedBy   uint      `gorm:"index" json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Name      string         `gorm:"type:varchar(100)" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Incident represents a downtime event for a source
// Phase 4: Incident Management
type Incident struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	SourceID    string     `gorm:"type:uuid;not null;index" json:"source_id"`
	Status      string     `gorm:"type:varchar(20);default:'OPEN';index" json:"status"` // OPEN, RESOLVED
	StartedAt   time.Time  `gorm:"not null;index" json:"started_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
	DurationSec int64      `json:"duration_sec"`
	Message     string     `gorm:"type:text" json:"message"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// APMThreshold stores latency thresholds for alerting
// Phase 3: APM Threshold
type APMThreshold struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SourceID    string    `gorm:"type:uuid;not null;index" json:"source_id"`
	Endpoint    string    `gorm:"type:varchar(255);not null" json:"endpoint"`
	P95Limit    int       `gorm:"not null;default:1000" json:"p95_limit"` // milliseconds
	P99Limit    int       `gorm:"not null;default:2000" json:"p99_limit"` // milliseconds
	EmailNotify bool      `gorm:"not null;default:true" json:"email_notify"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Session tracks user login sessions
// Phase 2: Activity Monitor
type Session struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       string    `gorm:"type:varchar(100);not null;index" json:"user_id"`
	SourceID     string    `gorm:"type:uuid;not null;index" json:"source_id"`
	AuthMethod   string    `gorm:"type:varchar(50);not null" json:"auth_method"`
	IPAddress    string    `gorm:"type:varchar(45)" json:"ip_address"`
	Browser      string    `gorm:"type:varchar(100)" json:"browser"`
	Device       string    `gorm:"type:varchar(100)" json:"device"`
	IsActive     bool      `gorm:"not null;default:true" json:"is_active"`
	LastActivity time.Time `json:"last_activity"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
