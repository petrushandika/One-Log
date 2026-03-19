package repository

import (
	"time"
)

func (r *logRepository) GetSlowQueryTrend(sourceID string, days int) ([]map[string]interface{}, error) {
	fromDate := time.Now().AddDate(0, 0, -days)

	sql := `
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as count,
			AVG((context->>'duration_ms')::int) as avg_duration
		FROM log_entries
		WHERE category = 'PERFORMANCE'
			AND created_at >= ?
			AND jsonb_exists(context, 'duration_ms')
			AND (context->>'duration_ms')::int > 2000
	`
	args := []interface{}{fromDate}

	if sourceID != "" {
		sql += " AND source_id = ?"
		args = append(args, sourceID)
	}

	sql += " GROUP BY DATE(created_at) ORDER BY date"

	rows, err := r.db.Raw(sql, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []map[string]interface{}{}
	for rows.Next() {
		var date string
		var count int
		var avgDuration float64

		if err := rows.Scan(&date, &count, &avgDuration); err != nil {
			continue
		}
		out = append(out, map[string]interface{}{
			"date":         date,
			"count":        count,
			"avg_duration": avgDuration,
		})
	}
	return out, nil
}

func (r *logRepository) CalculateApdexScore(sourceID string, endpoint string, thresholdMs int) (*ApdexResult, error) {
	sql := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN (context->>'duration_ms')::int <= ? THEN 1 END) as satisfied,
			COUNT(CASE WHEN (context->>'duration_ms')::int > ? AND (context->>'duration_ms')::int <= ? * 4 THEN 1 END) as tolerating,
			COUNT(CASE WHEN (context->>'duration_ms')::int > ? * 4 THEN 1 END) as frustrated
		FROM log_entries
		WHERE category = 'PERFORMANCE'
			AND jsonb_exists(context, 'duration_ms')
	`
	args := []interface{}{thresholdMs, thresholdMs, thresholdMs, thresholdMs}

	if sourceID != "" {
		sql += " AND source_id = ?"
		args = append(args, sourceID)
	}

	if endpoint != "" {
		sql += " AND context->>'endpoint' = ?"
		args = append(args, endpoint)
	}

	var total, satisfied, tolerating, frustrated int
	err := r.db.Raw(sql, args...).Row().Scan(&total, &satisfied, &tolerating, &frustrated)
	if err != nil {
		return nil, err
	}

	// Calculate Apdex score
	var score float64
	if total > 0 {
		score = (float64(satisfied) + float64(tolerating)/2.0) / float64(total)
	}

	return &ApdexResult{
		Score:      score,
		Satisfied:  satisfied,
		Tolerating: tolerating,
		Frustrated: frustrated,
		Total:      total,
	}, nil
}

func (r *logRepository) GetEndpointLatencyStats(sourceID string, endpoint string) (map[string]interface{}, error) {
	sql := `
		SELECT 
			PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY (context->>'duration_ms')::int) as p50,
			PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY (context->>'duration_ms')::int) as p95,
			PERCENTILE_CONT(0.99) WITHIN GROUP (ORDER BY (context->>'duration_ms')::int) as p99,
			COUNT(*) as count
		FROM log_entries
		WHERE category = 'PERFORMANCE'
			AND jsonb_exists(context, 'duration_ms')
	`
	args := []interface{}{}

	if sourceID != "" {
		sql += " AND source_id = ?"
		args = append(args, sourceID)
	}

	if endpoint != "" {
		sql += " AND context->>'endpoint' = ?"
		args = append(args, endpoint)
	}

	var p50, p95, p99 float64
	var count int64
	err := r.db.Raw(sql, args...).Row().Scan(&p50, &p95, &p99, &count)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"p50_ms": p50,
		"p95_ms": p95,
		"p99_ms": p99,
		"count":  count,
	}, nil
}
