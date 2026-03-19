package repository

import (
	"time"

	"gorm.io/gorm"
)

// ActivityAnalyticsRepository defines analytics methods for activity logs
type ActivityAnalyticsRepository interface {
	GetAuthMethodBreakdown(sourceID string, days int) (map[string]int64, error)
	GetLoginTimeline(sourceID string, days int) ([]map[string]interface{}, error)
	GetFailedLoginHeatmap(sourceID string, days int) ([]map[string]interface{}, error)
	GetRecentSessions(limit int, offset int, sourceID string) ([]map[string]interface{}, int64, error)
}

type activityAnalyticsRepository struct {
	db *gorm.DB
}

// NewActivityAnalyticsRepository creates a new activity analytics repository
func NewActivityAnalyticsRepository(db *gorm.DB) ActivityAnalyticsRepository {
	return &activityAnalyticsRepository{db: db}
}

func (r *activityAnalyticsRepository) GetAuthMethodBreakdown(sourceID string, days int) (map[string]int64, error) {
	sql := `
SELECT 
    context->>'auth_method' as auth_method,
    count(*) as count
FROM log_entries
WHERE category = 'AUTH_EVENT'
AND created_at >= NOW() - INTERVAL '1 day' * ?
`
	args := []interface{}{days}

	if sourceID != "" {
		sql += " AND source_id = ?"
		args = append(args, sourceID)
	}

	sql += " GROUP BY context->>'auth_method'"

	rows, err := r.db.Raw(sql, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var method string
		var count int64
		if err := rows.Scan(&method, &count); err != nil {
			continue
		}
		result[method] = count
	}
	return result, nil
}

func (r *activityAnalyticsRepository) GetLoginTimeline(sourceID string, days int) ([]map[string]interface{}, error) {
	sql := `
SELECT 
    DATE(created_at) as date,
    count(*) FILTER (WHERE context->>'success' = 'true') as login_success,
    count(*) FILTER (WHERE context->>'success' = 'false') as login_failed
FROM log_entries
WHERE category = 'AUTH_EVENT'
AND created_at >= NOW() - INTERVAL '1 day' * ?
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
		var success, failed int64
		if err := rows.Scan(&date, &success, &failed); err != nil {
			continue
		}
		out = append(out, map[string]interface{}{
			"date":          date.Format("2006-01-02"),
			"login_success": success,
			"login_failed":  failed,
		})
	}
	return out, nil
}

func (r *activityAnalyticsRepository) GetFailedLoginHeatmap(sourceID string, days int) ([]map[string]interface{}, error) {
	sql := `
SELECT 
    extract(dow from created_at) as day_of_week,
    extract(hour from created_at) as hour_of_day,
    count(*) as failed_count
FROM log_entries
WHERE category = 'AUTH_EVENT'
AND context->>'success' = 'false'
AND created_at >= NOW() - INTERVAL '1 day' * ?
`
	args := []interface{}{days}

	if sourceID != "" {
		sql += " AND source_id = ?"
		args = append(args, sourceID)
	}

	sql += " GROUP BY day_of_week, hour_of_day ORDER BY day_of_week, hour_of_day"

	rows, err := r.db.Raw(sql, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []map[string]interface{}{}
	for rows.Next() {
		var dayOfWeek, hourOfDay, count int64
		if err := rows.Scan(&dayOfWeek, &hourOfDay, &count); err != nil {
			continue
		}
		out = append(out, map[string]interface{}{
			"day_of_week":  dayOfWeek,
			"hour_of_day":  hourOfDay,
			"failed_count": count,
		})
	}
	return out, nil
}

func (r *activityAnalyticsRepository) GetRecentSessions(limit int, offset int, sourceID string) ([]map[string]interface{}, int64, error) {
	// Check if sessions table exists
	var tableExists bool
	err := r.db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'sessions')").Scan(&tableExists).Error
	if err != nil || !tableExists {
		// Return empty result if table doesn't exist
		return []map[string]interface{}{}, 0, nil
	}

	sql := `
SELECT 
    id,
    user_id,
    source_id,
    auth_method,
    ip_address,
    browser,
    device,
    created_at,
    last_activity
FROM sessions
WHERE is_active = true
`
	args := []interface{}{}

	if sourceID != "" {
		sql += " AND source_id = ?"
		args = append(args, sourceID)
	}

	// Count total
	countSQL := "SELECT count(*) FROM sessions WHERE is_active = true"
	if sourceID != "" {
		countSQL += " AND source_id = ?"
	}

	var total int64
	if err := r.db.Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	sql += " ORDER BY last_activity DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Raw(sql, args...).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := []map[string]interface{}{}
	for rows.Next() {
		var id int64
		var userID, authMethod, ipAddress string
		var sourceIdStr string
		var browser, device string
		var createdAt, lastActivity time.Time

		if err := rows.Scan(&id, &userID, &sourceIdStr, &authMethod, &ipAddress, &browser, &device, &createdAt, &lastActivity); err != nil {
			continue
		}
		out = append(out, map[string]interface{}{
			"id":            id,
			"user_id":       userID,
			"source_id":     sourceIdStr,
			"auth_method":   authMethod,
			"ip_address":    ipAddress,
			"browser":       browser,
			"device":        device,
			"created_at":    createdAt,
			"last_activity": lastActivity,
		})
	}
	return out, total, nil
}
