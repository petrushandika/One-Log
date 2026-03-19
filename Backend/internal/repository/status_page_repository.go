package repository

import (
	"time"

	"github.com/petrushandika/one-log/internal/domain"
	"gorm.io/gorm"
)

type StatusPageRepository interface {
	// Status Page Config
	CreateStatusPageConfig(config *domain.StatusPageConfig) error
	GetStatusPageConfig(sourceID string) (*domain.StatusPageConfig, error)
	GetStatusPageConfigBySlug(slug string) (*domain.StatusPageConfig, error)
	UpdateStatusPageConfig(config *domain.StatusPageConfig) error
	DeleteStatusPageConfig(sourceID string) error
	ListStatusPageConfigs(page, limit int) ([]domain.StatusPageConfig, int64, error)

	// Uptime Stats
	GetUptimeStats(sourceID string, days int) (*domain.UptimeStats, error)
	GetAllUptimeStats(days int) ([]domain.UptimeStats, error)

	// Embed Widget
	CreateEmbedToken(sourceID string) (*domain.StatusPageEmbed, error)
	GetEmbedByToken(token string) (*domain.StatusPageEmbed, error)
	GetEmbedBySource(sourceID string) (*domain.StatusPageEmbed, error)
	DeleteEmbedToken(token string) error
}

type statusPageRepository struct {
	db *gorm.DB
}

func NewStatusPageRepository(db *gorm.DB) StatusPageRepository {
	return &statusPageRepository{db: db}
}

func (r *statusPageRepository) CreateStatusPageConfig(config *domain.StatusPageConfig) error {
	return r.db.Create(config).Error
}

func (r *statusPageRepository) GetStatusPageConfig(sourceID string) (*domain.StatusPageConfig, error) {
	var config domain.StatusPageConfig
	err := r.db.Where("source_id = ?", sourceID).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *statusPageRepository) GetStatusPageConfigBySlug(slug string) (*domain.StatusPageConfig, error) {
	var config domain.StatusPageConfig
	err := r.db.Where("slug = ?", slug).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *statusPageRepository) UpdateStatusPageConfig(config *domain.StatusPageConfig) error {
	return r.db.Save(config).Error
}

func (r *statusPageRepository) DeleteStatusPageConfig(sourceID string) error {
	return r.db.Where("source_id = ?", sourceID).Delete(&domain.StatusPageConfig{}).Error
}

func (r *statusPageRepository) ListStatusPageConfigs(page, limit int) ([]domain.StatusPageConfig, int64, error) {
	var configs []domain.StatusPageConfig
	var total int64

	r.db.Model(&domain.StatusPageConfig{}).Count(&total)

	offset := (page - 1) * limit
	err := r.db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&configs).Error

	return configs, total, err
}

func (r *statusPageRepository) GetUptimeStats(sourceID string, days int) (*domain.UptimeStats, error) {
	fromDate := time.Now().AddDate(0, 0, -days)

	// Get source info
	var source domain.Source
	if err := r.db.First(&source, "id = ?", sourceID).Error; err != nil {
		return nil, err
	}

	// Calculate total downtime and incident count
	var totalDowntime int
	var incidentCount int64

	r.db.Model(&domain.Incident{}).
		Where("source_id = ? AND started_at >= ?", sourceID, fromDate).
		Count(&incidentCount)

	r.db.Raw(`
		SELECT COALESCE(SUM(EXTRACT(EPOCH FROM (resolved_at - started_at))), 0) as downtime
		FROM incidents
		WHERE source_id = ? AND started_at >= ? AND resolved_at IS NOT NULL
	`, sourceID, fromDate).Scan(&totalDowntime)

	// Calculate uptime percentage
	totalSeconds := days * 24 * 60 * 60
	uptimePercent := 100.0
	if totalSeconds > 0 {
		uptimePercent = 100.0 - (float64(totalDowntime) / float64(totalSeconds) * 100)
	}

	// Get last incident
	var lastIncident domain.Incident
	lastIncidentAt := r.db.Where("source_id = ?", sourceID).Order("started_at DESC").First(&lastIncident)

	// Get last log time for this source
	var lastLog domain.LogEntry
	lastLogAt := time.Now()
	r.db.Where("source_id = ?", sourceID).Order("created_at DESC").First(&lastLog)
	if lastLog.ID > 0 {
		lastLogAt = lastLog.CreatedAt
	}

	stats := &domain.UptimeStats{
		SourceID:         sourceID,
		SourceName:       source.Name,
		CurrentStatus:    source.Status,
		UptimePercent30d: uptimePercent,
		TotalDowntime30d: totalDowntime,
		IncidentCount30d: int(incidentCount),
		LastCheckedAt:    lastLogAt,
	}

	if lastIncidentAt.Error == nil {
		stats.LastIncidentAt = &lastIncident.StartedAt
	}

	return stats, nil
}

func (r *statusPageRepository) GetAllUptimeStats(days int) ([]domain.UptimeStats, error) {
	var sources []domain.Source
	if err := r.db.Find(&sources).Error; err != nil {
		return nil, err
	}

	var stats []domain.UptimeStats
	for _, source := range sources {
		stat, err := r.GetUptimeStats(source.ID, days)
		if err != nil {
			continue
		}
		stats = append(stats, *stat)
	}

	return stats, nil
}

func (r *statusPageRepository) CreateEmbedToken(sourceID string) (*domain.StatusPageEmbed, error) {
	// Generate unique token
	token := generateUniqueToken()

	embed := &domain.StatusPageEmbed{
		SourceID:   sourceID,
		EmbedToken: token,
		Theme:      "auto",
		CreatedAt:  time.Now(),
	}

	err := r.db.Create(embed).Error
	return embed, err
}

func (r *statusPageRepository) GetEmbedByToken(token string) (*domain.StatusPageEmbed, error) {
	var embed domain.StatusPageEmbed
	err := r.db.Where("embed_token = ?", token).First(&embed).Error
	if err != nil {
		return nil, err
	}
	return &embed, nil
}

func (r *statusPageRepository) GetEmbedBySource(sourceID string) (*domain.StatusPageEmbed, error) {
	var embed domain.StatusPageEmbed
	err := r.db.Where("source_id = ?", sourceID).First(&embed).Error
	if err != nil {
		return nil, err
	}
	return &embed, nil
}

func (r *statusPageRepository) DeleteEmbedToken(token string) error {
	return r.db.Where("embed_token = ?", token).Delete(&domain.StatusPageEmbed{}).Error
}

func generateUniqueToken() string {
	// Simple token generation - in production use crypto/rand
	return "embed_" + time.Now().Format("20060102150405") + "_" + generateRandomString(16)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}
