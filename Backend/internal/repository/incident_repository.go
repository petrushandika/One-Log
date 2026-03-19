package repository

import (
	"time"

	"github.com/petrushandika/one-log/internal/domain"
	"gorm.io/gorm"
)

// IncidentRepository defines data access methods for incidents
type IncidentRepository interface {
	Create(incident *domain.Incident) error
	FindOpenBySource(sourceID string) (*domain.Incident, error)
	Resolve(id uint, message string) error
	List(limit int, offset int, sourceID string, status string) ([]domain.Incident, int64, error)
	GetTimeline(sourceID string, days int) ([]map[string]interface{}, error)
}

type incidentRepository struct {
	db *gorm.DB
}

// NewIncidentRepository creates a new incident repository
func NewIncidentRepository(db *gorm.DB) IncidentRepository {
	return &incidentRepository{db: db}
}

func (r *incidentRepository) Create(incident *domain.Incident) error {
	return r.db.Create(incident).Error
}

func (r *incidentRepository) FindOpenBySource(sourceID string) (*domain.Incident, error) {
	var incident domain.Incident
	err := r.db.Where("source_id = ? AND status = 'OPEN'", sourceID).First(&incident).Error
	if err != nil {
		return nil, err
	}
	return &incident, nil
}

func (r *incidentRepository) Resolve(id uint, message string) error {
	now := time.Now()
	var incident domain.Incident
	if err := r.db.First(&incident, id).Error; err != nil {
		return err
	}

	duration := int64(now.Sub(incident.StartedAt).Seconds())

	return r.db.Model(&incident).Updates(map[string]interface{}{
		"status":       "RESOLVED",
		"resolved_at":  now,
		"duration_sec": duration,
		"message":      gorm.Expr("CASE WHEN message = '' THEN ? ELSE message || '\n' || ? END", message, message),
		"updated_at":   now,
	}).Error
}

func (r *incidentRepository) List(limit int, offset int, sourceID string, status string) ([]domain.Incident, int64, error) {
	var incidents []domain.Incident
	var total int64

	query := r.db.Model(&domain.Incident{})

	if sourceID != "" {
		query = query.Where("source_id = ?", sourceID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("started_at desc").Limit(limit).Offset(offset).Find(&incidents).Error; err != nil {
		return nil, 0, err
	}

	return incidents, total, nil
}

func (r *incidentRepository) GetTimeline(sourceID string, days int) ([]map[string]interface{}, error) {
	sql := `
SELECT 
    DATE(created_at) as date,
    count(*) filter (where status = 'OPEN') as open_count,
    count(*) filter (where status = 'RESOLVED') as resolved_count,
    sum(duration_sec) filter (where status = 'RESOLVED') as total_downtime_sec
FROM incidents
WHERE created_at >= NOW() - INTERVAL '1 day' * ?
`
	args := []interface{}{days}

	if sourceID != "" {
		sql += " AND source_id = ?"
		args = append(args, sourceID)
	}

	sql += " GROUP BY DATE(created_at) ORDER BY date ASC"

	rows, err := r.db.Raw(sql, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []map[string]interface{}{}
	for rows.Next() {
		var date time.Time
		var openCount, resolvedCount int64
		var totalDowntime *int64
		if err := rows.Scan(&date, &openCount, &resolvedCount, &totalDowntime); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"date":               date.Format("2006-01-02"),
			"open_count":         openCount,
			"resolved_count":     resolvedCount,
			"total_downtime_sec": totalDowntime,
		})
	}
	return out, nil
}
