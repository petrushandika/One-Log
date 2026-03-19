package domain

import "time"

// ActivityFeed represents a user activity entry for the activity feed
type ActivityFeed struct {
	ID           uint                   `json:"id" gorm:"primaryKey"`
	UserID       string                 `json:"user_id" gorm:"index"`
	SourceID     string                 `json:"source_id" gorm:"index"`
	Action       string                 `json:"action"` // page_view, create, update, delete, export
	ResourceType string                 `json:"resource_type"`
	ResourceID   string                 `json:"resource_id"`
	Context      map[string]interface{} `json:"context" gorm:"serializer:json"`
	IPAddress    string                 `json:"ip_address"`
	CreatedAt    time.Time              `json:"created_at"`
}

// ActiveUser represents a top active user statistics
type ActiveUser struct {
	UserID        string    `json:"user_id"`
	SourceID      string    `json:"source_id"`
	ActivityCount int       `json:"activity_count"`
	LastActive    time.Time `json:"last_active"`
}

// ResourceActivity represents activity count for a specific resource
type ResourceActivity struct {
	ResourceType string    `json:"resource_type"`
	ResourceID   string    `json:"resource_id"`
	Action       string    `json:"action"`
	Count        int       `json:"count"`
	LastAccessed time.Time `json:"last_accessed"`
}

// BeforeAfterDiff stores before/after values for audit trail
type BeforeAfterDiff struct {
	Before map[string]interface{} `json:"before" gorm:"serializer:json"`
	After  map[string]interface{} `json:"after" gorm:"serializer:json"`
	Diff   map[string]interface{} `json:"diff" gorm:"serializer:json"` // Computed differences
}

// ComplianceExport represents an audit trail export record
type ComplianceExport struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	SourceID    string     `json:"source_id"`
	Format      string     `json:"format"` // PDF, CSV
	DateFrom    time.Time  `json:"date_from"`
	DateTo      time.Time  `json:"date_to"`
	Status      string     `json:"status"` // pending, processing, completed, failed
	FileURL     string     `json:"file_url"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}
