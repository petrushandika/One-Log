package domain

import "time"

// StatusPageConfig stores configuration for public status pages
type StatusPageConfig struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	SourceID    string    `json:"source_id" gorm:"uniqueIndex"`
	Slug        string    `json:"slug" gorm:"uniqueIndex"` // Custom URL slug
	Title       string    `json:"title"`
	Description string    `json:"description"`
	LogoURL     string    `json:"logo_url"`
	IsPublic    bool      `json:"is_public"` // Whether accessible without auth
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UptimeStats holds uptime statistics for display
type UptimeStats struct {
	SourceID         string     `json:"source_id"`
	SourceName       string     `json:"source_name"`
	CurrentStatus    string     `json:"current_status"`
	UptimePercent30d float64    `json:"uptime_percent_30d"`
	UptimePercent90d float64    `json:"uptime_percent_90d"`
	TotalDowntime30d int        `json:"total_downtime_30d"` // in seconds
	IncidentCount30d int        `json:"incident_count_30d"`
	LastCheckedAt    time.Time  `json:"last_checked_at"`
	LastIncidentAt   *time.Time `json:"last_incident_at"`
}

// StatusPageEmbed holds embed widget configuration
type StatusPageEmbed struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	SourceID   string    `json:"source_id"`
	EmbedToken string    `json:"embed_token" gorm:"uniqueIndex"`
	Theme      string    `json:"theme"` // light, dark, auto
	CreatedAt  time.Time `json:"created_at"`
}
