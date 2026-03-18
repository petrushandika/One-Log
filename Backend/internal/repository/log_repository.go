package repository

import (
	"fmt"
	"time"

	"github.com/petrushandika/one-log/internal/domain"
	"gorm.io/gorm"
)

// Definition LogRepository interface
type LogRepository interface {
	Create(log *domain.LogEntry) error
	FindAll(limit int, offset int, sourceID, level, category string, userID uint, from, to *time.Time) ([]domain.LogEntry, int64, error)
	FindByID(id uint) (*domain.LogEntry, error)
	Update(log *domain.LogEntry) error
	DeleteOlderThan(days int) error
	GetStatsOverview(userID uint) (map[string]interface{}, error)
	CountFailedAttempts(ip string, durationMinutes int) (int64, error)
	GetActivitySummary(userID uint) (map[string]interface{}, error)

	FindActivity(limit int, offset int, sourceID string, categories []string, eventType, authMethod, subjectUserID string, from, to *time.Time, ownerUserID uint) ([]domain.LogEntry, int64, error)
	GetActivitySummaryByPeriod(period time.Duration, sourceID string, ownerUserID uint) (map[string]interface{}, error)
	GetUserActivity(userID string, period time.Duration, categories []string, ownerUserID uint) (map[string]interface{}, error)
	FindSuspiciousActivity(limit int, offset int, period time.Duration, sourceID string, ownerUserID uint) ([]domain.LogEntry, int64, error)

	GetAPMEndpointStats(period time.Duration, sourceID string, ownerUserID uint) ([]map[string]interface{}, error)

	UpsertIssueFromLog(log *domain.LogEntry) error
	ListIssues(limit int, offset int, sourceID, status string, ownerUserID uint) ([]domain.Issue, int64, error)
	GetIssueByFingerprint(fingerprint string, ownerUserID uint) (*domain.Issue, error)
	UpdateIssueStatus(fingerprint string, status string, ownerUserID uint) (*domain.Issue, error)
	ListIssueLogs(limit int, offset int, fingerprint string, ownerUserID uint) ([]domain.LogEntry, int64, error)
}

// Struct private for implementation
type logRepository struct {
	db *gorm.DB
}

// Constructor
func NewLogRepository(db *gorm.DB) LogRepository {
	return &logRepository{db: db}
}

// Implementation method Create
func (r *logRepository) Create(log *domain.LogEntry) error {
	return r.db.Create(log).Error
}

func (r *logRepository) FindAll(limit int, offset int, sourceID, level, category string, userID uint, from, to *time.Time) ([]domain.LogEntry, int64, error) {
	var logs []domain.LogEntry
	var total int64
	query := r.db.Model(&domain.LogEntry{})

	if userID > 0 {
		query = query.Joins("JOIN sources ON sources.id = log_entries.source_id").Where("sources.user_id = ?", userID)
	}

	if sourceID != "" {
		query = query.Where("log_entries.source_id = ?", sourceID)
	}
	if level != "" {
		query = query.Where("log_entries.level = ?", level)
	}
	if category != "" {
		query = query.Where("log_entries.category = ?", category)
	}
	if from != nil {
		query = query.Where("log_entries.created_at >= ?", *from)
	}
	if to != nil {
		query = query.Where("log_entries.created_at <= ?", *to)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("log_entries.created_at desc").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}

func (r *logRepository) FindByID(id uint) (*domain.LogEntry, error) {
	var log domain.LogEntry
	err := r.db.First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *logRepository) Update(log *domain.LogEntry) error {
	return r.db.Save(log).Error
}

func (r *logRepository) DeleteOlderThan(days int) error {
	// Fase 2 Immutable Logs Guard:
	// Logs categorized as AUDIT_TRAIL should NOT be deleted by general retention policies.
	return r.db.Where("created_at < NOW() - INTERVAL '1 day' * ? AND category != 'AUDIT_TRAIL'", days).Delete(&domain.LogEntry{}).Error
}

func (r *logRepository) GetStatsOverview(userID uint) (map[string]interface{}, error) {
	queryTotal := r.db.Model(&domain.LogEntry{})
	queryBreakdown := r.db.Model(&domain.LogEntry{})

	if userID > 0 {
		queryTotal = queryTotal.Joins("JOIN sources ON sources.id = log_entries.source_id").Where("sources.user_id = ?", userID)
		queryBreakdown = queryBreakdown.Joins("JOIN sources ON sources.id = log_entries.source_id").Where("sources.user_id = ?", userID)
	}

	var total int64
	if err := queryTotal.Count(&total).Error; err != nil {
		return nil, err
	}

	type result struct {
		Level string
		Count int64
	}
	var bResult []result
	if err := queryBreakdown.Select("log_entries.level, count(*) as count").Group("log_entries.level").Find(&bResult).Error; err != nil {
		return nil, err
	}

	stats := map[string]interface{}{"total": total}
	for _, row := range bResult {
		stats[row.Level] = row.Count
	}

	return stats, nil
}

func (r *logRepository) CountFailedAttempts(ip string, durationMinutes int) (int64, error) {
	var count int64
	err := r.db.Model(&domain.LogEntry{}).
		Where("ip_address = ? AND category = 'AUTH_EVENT' AND level = 'WARN' AND created_at >= NOW() - INTERVAL '1 minute' * ?", ip, durationMinutes).
		Count(&count).Error
	return count, err
}

func (r *logRepository) GetActivitySummary(userID uint) (map[string]interface{}, error) {
	queryMethod := r.db.Model(&domain.LogEntry{}).Where("log_entries.category = 'AUTH_EVENT'")
	queryType := r.db.Model(&domain.LogEntry{}).Where("log_entries.category = 'AUTH_EVENT'")

	if userID > 0 {
		queryMethod = queryMethod.Joins("JOIN sources ON sources.id = log_entries.source_id").Where("sources.user_id = ?", userID)
		queryType = queryType.Joins("JOIN sources ON sources.id = log_entries.source_id").Where("sources.user_id = ?", userID)
	}

	type methodResult struct {
		AuthMethod string `gorm:"column:auth_method"`
		Count      int64  `gorm:"column:count"`
	}
	type typeResult struct {
		EventType string `gorm:"column:event_type"`
		Count     int64  `gorm:"column:count"`
	}

	var methods []methodResult
	var types []typeResult

	// Aggregate by Auth Method
	err := queryMethod.Select("log_entries.context->>'auth_method' as auth_method, count(*) as count").
		Group("log_entries.context->>'auth_method'").
		Scan(&methods).Error
	if err != nil {
		return nil, err
	}

	// Aggregate by Event Type
	err = queryType.Select("log_entries.context->>'event_type' as event_type, count(*) as count").
		Group("log_entries.context->>'event_type'").
		Scan(&types).Error
	if err != nil {
		return nil, err
	}

	methodMap := make(map[string]int64)
	for _, m := range methods {
		if m.AuthMethod == "" {
			methodMap["unknown"] = m.Count
		} else {
			methodMap[m.AuthMethod] = m.Count
		}
	}

	typeMap := make(map[string]int64)
	for _, t := range types {
		if t.EventType == "" {
			typeMap["unknown"] = t.Count
		} else {
			typeMap[t.EventType] = t.Count
		}
	}

	return map[string]interface{}{
		"auth_methods": methodMap,
		"event_types":  typeMap,
	}, nil
}

func (r *logRepository) FindActivity(limit int, offset int, sourceID string, categories []string, eventType, authMethod, subjectUserID string, from, to *time.Time, ownerUserID uint) ([]domain.LogEntry, int64, error) {
	var logs []domain.LogEntry
	var total int64

	query := r.db.Model(&domain.LogEntry{})
	if ownerUserID > 0 {
		query = query.Joins("JOIN sources ON sources.id = log_entries.source_id").Where("sources.user_id = ?", ownerUserID)
	}
	if sourceID != "" {
		query = query.Where("log_entries.source_id = ?", sourceID)
	}
	if len(categories) > 0 {
		query = query.Where("log_entries.category IN ?", categories)
	}
	if eventType != "" {
		query = query.Where("log_entries.context->>'event_type' = ?", eventType)
	}
	if authMethod != "" {
		query = query.Where("log_entries.context->>'auth_method' = ?", authMethod)
	}
	if subjectUserID != "" {
		query = query.Where("log_entries.context->>'user_id' = ?", subjectUserID)
	}
	if from != nil {
		query = query.Where("log_entries.created_at >= ?", *from)
	}
	if to != nil {
		query = query.Where("log_entries.created_at <= ?", *to)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("log_entries.created_at desc").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}

func (r *logRepository) GetActivitySummaryByPeriod(period time.Duration, sourceID string, ownerUserID uint) (map[string]interface{}, error) {
	since := time.Now().Add(-period)
	queryMethod := r.db.Model(&domain.LogEntry{}).Where("log_entries.category = 'AUTH_EVENT' AND log_entries.created_at >= ?", since)
	queryType := r.db.Model(&domain.LogEntry{}).Where("log_entries.category = 'AUTH_EVENT' AND log_entries.created_at >= ?", since)
	querySuspicious := r.db.Model(&domain.LogEntry{}).Where("log_entries.created_at >= ?", since)

	if ownerUserID > 0 {
		queryMethod = queryMethod.Joins("JOIN sources ON sources.id = log_entries.source_id").Where("sources.user_id = ?", ownerUserID)
		queryType = queryType.Joins("JOIN sources ON sources.id = log_entries.source_id").Where("sources.user_id = ?", ownerUserID)
		querySuspicious = querySuspicious.Joins("JOIN sources ON sources.id = log_entries.source_id").Where("sources.user_id = ?", ownerUserID)
	}
	if sourceID != "" {
		queryMethod = queryMethod.Where("log_entries.source_id = ?", sourceID)
		queryType = queryType.Where("log_entries.source_id = ?", sourceID)
		querySuspicious = querySuspicious.Where("log_entries.source_id = ?", sourceID)
	}

	type methodResult struct {
		AuthMethod string `gorm:"column:auth_method"`
		Count      int64  `gorm:"column:count"`
	}
	type typeResult struct {
		EventType string `gorm:"column:event_type"`
		Count     int64  `gorm:"column:count"`
	}
	type trendRow struct {
		Day   string `gorm:"column:day"`
		Count int64  `gorm:"column:count"`
	}

	var methods []methodResult
	var types []typeResult
	if err := queryMethod.Select("log_entries.context->>'auth_method' as auth_method, count(*) as count").
		Group("log_entries.context->>'auth_method'").
		Scan(&methods).Error; err != nil {
		return nil, err
	}
	if err := queryType.Select("log_entries.context->>'event_type' as event_type, count(*) as count").
		Group("log_entries.context->>'event_type'").
		Scan(&types).Error; err != nil {
		return nil, err
	}

	methodMap := make(map[string]int64)
	for _, m := range methods {
		if m.AuthMethod == "" {
			methodMap["unknown"] += m.Count
		} else {
			methodMap[m.AuthMethod] = m.Count
		}
	}
	typeMap := make(map[string]int64)
	for _, t := range types {
		if t.EventType == "" {
			typeMap["unknown"] += t.Count
		} else {
			typeMap[t.EventType] = t.Count
		}
	}

	// Failed login trend (daily)
	var failedTrend []trendRow
	_ = querySuspicious.
		Where("log_entries.category = 'AUTH_EVENT' AND log_entries.context->>'event_type' = 'login_failed'").
		Select("to_char(date_trunc('day', log_entries.created_at), 'YYYY-MM-DD') as day, count(*) as count").
		Group("day").
		Order("day asc").
		Scan(&failedTrend).Error

	// Suspicious count
	var suspiciousCount int64
	_ = querySuspicious.
		Where("(log_entries.category = 'AUTH_EVENT' AND (log_entries.context->>'event_type' = 'suspicious_login' OR log_entries.level IN ('ERROR','CRITICAL'))) OR log_entries.category = 'SECURITY'").
		Count(&suspiciousCount).Error

	return map[string]interface{}{
		"since":              since.UTC().Format(time.RFC3339),
		"by_auth_method":     methodMap,
		"by_event_type":      typeMap,
		"failed_login_trend": failedTrend,
		"suspicious_count":   suspiciousCount,
	}, nil
}

func (r *logRepository) GetUserActivity(userID string, period time.Duration, categories []string, ownerUserID uint) (map[string]interface{}, error) {
	since := time.Now().Add(-period)
	query := r.db.Model(&domain.LogEntry{}).Where("log_entries.created_at >= ? AND log_entries.context->>'user_id' = ?", since, userID)
	if ownerUserID > 0 {
		query = query.Joins("JOIN sources ON sources.id = log_entries.source_id").Where("sources.user_id = ?", ownerUserID)
	}
	if len(categories) > 0 {
		query = query.Where("log_entries.category IN ?", categories)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var recent []domain.LogEntry
	if err := query.Order("log_entries.created_at desc").Limit(50).Find(&recent).Error; err != nil {
		return nil, err
	}

	// Distinct apps accessed
	type appRow struct{ SourceID string }
	var apps []appRow
	_ = query.Select("distinct log_entries.source_id as source_id").Scan(&apps).Error
	appIDs := make([]string, 0, len(apps))
	for _, a := range apps {
		appIDs = append(appIDs, a.SourceID)
	}

	return map[string]interface{}{
		"user_id":       userID,
		"since":         since.UTC().Format(time.RFC3339),
		"total_events":  total,
		"apps_accessed": appIDs,
		"recent_events": recent,
	}, nil
}

func (r *logRepository) FindSuspiciousActivity(limit int, offset int, period time.Duration, sourceID string, ownerUserID uint) ([]domain.LogEntry, int64, error) {
	var logs []domain.LogEntry
	var total int64
	since := time.Now().Add(-period)

	query := r.db.Model(&domain.LogEntry{}).
		Where("log_entries.created_at >= ?", since).
		Where("(log_entries.category = 'AUTH_EVENT' AND (log_entries.context->>'event_type' = 'suspicious_login' OR log_entries.level IN ('ERROR','CRITICAL'))) OR log_entries.category IN ('SECURITY','AUDIT_TRAIL')")

	if ownerUserID > 0 {
		query = query.Joins("JOIN sources ON sources.id = log_entries.source_id").Where("sources.user_id = ?", ownerUserID)
	}
	if sourceID != "" {
		query = query.Where("log_entries.source_id = ?", sourceID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("log_entries.created_at desc").Limit(limit).Offset(offset).Find(&logs).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query suspicious activity: %w", err)
	}
	return logs, total, nil
}

func (r *logRepository) GetAPMEndpointStats(period time.Duration, sourceID string, ownerUserID uint) ([]map[string]interface{}, error) {
	since := time.Now().Add(-period)

	// Build base SQL with optional ownership/source filters.
	sql := `
SELECT
  COALESCE(log_entries.context->>'endpoint', 'unknown') AS endpoint,
  COUNT(*) AS count,
  percentile_cont(0.50) WITHIN GROUP (ORDER BY (log_entries.context->>'duration_ms')::numeric) AS p50,
  percentile_cont(0.95) WITHIN GROUP (ORDER BY (log_entries.context->>'duration_ms')::numeric) AS p95,
  percentile_cont(0.99) WITHIN GROUP (ORDER BY (log_entries.context->>'duration_ms')::numeric) AS p99
FROM log_entries
`
	args := []interface{}{}

	joins := ""
	where := "WHERE log_entries.category = 'PERFORMANCE' AND log_entries.created_at >= ? AND jsonb_exists(log_entries.context, 'duration_ms')"
	args = append(args, since)

	if ownerUserID > 0 {
		joins += " JOIN sources ON sources.id = log_entries.source_id "
		where += " AND sources.user_id = ?"
		args = append(args, ownerUserID)
	}
	if sourceID != "" {
		where += " AND log_entries.source_id = ?"
		args = append(args, sourceID)
	}

	groupOrder := " GROUP BY endpoint ORDER BY p95 DESC"
	finalSQL := sql + joins + "\n" + where + "\n" + groupOrder

	rows, err := r.db.Raw(finalSQL, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []map[string]interface{}{}
	for rows.Next() {
		var endpoint string
		var count int64
		var p50, p95, p99 float64
		if err := rows.Scan(&endpoint, &count, &p50, &p95, &p99); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"endpoint": endpoint,
			"count":    count,
			"p50":      p50,
			"p95":      p95,
			"p99":      p99,
		})
	}
	return out, nil
}

func (r *logRepository) UpsertIssueFromLog(log *domain.LogEntry) error {
	if log.Fingerprint == "" {
		return nil
	}
	issue := domain.Issue{
		Fingerprint:     log.Fingerprint,
		SourceID:        log.SourceID,
		Status:          "OPEN",
		Category:        log.Category,
		Level:           log.Level,
		MessageSample:   log.Message,
		OccurrenceCount: 1,
		FirstSeenAt:     time.Now(),
		LastSeenAt:      time.Now(),
	}

	// PostgreSQL upsert
	return r.db.Exec(`
INSERT INTO issues (fingerprint, source_id, status, category, level, message_sample, occurrence_count, first_seen_at, last_seen_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, 1, ?, ?, NOW(), NOW())
ON CONFLICT (fingerprint)
DO UPDATE SET
  occurrence_count = issues.occurrence_count + 1,
  last_seen_at = EXCLUDED.last_seen_at,
  level = EXCLUDED.level,
  category = EXCLUDED.category,
  message_sample = EXCLUDED.message_sample,
  updated_at = NOW()
`, issue.Fingerprint, issue.SourceID, issue.Status, issue.Category, issue.Level, issue.MessageSample, issue.FirstSeenAt, issue.LastSeenAt).Error
}

func (r *logRepository) ListIssues(limit int, offset int, sourceID, status string, ownerUserID uint) ([]domain.Issue, int64, error) {
	var issues []domain.Issue
	var total int64
	query := r.db.Model(&domain.Issue{})
	if ownerUserID > 0 {
		query = query.Joins("JOIN sources ON sources.id = issues.source_id").Where("sources.user_id = ?", ownerUserID)
	}
	if sourceID != "" {
		query = query.Where("issues.source_id = ?", sourceID)
	}
	if status != "" {
		query = query.Where("issues.status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("issues.last_seen_at desc").Limit(limit).Offset(offset).Find(&issues).Error; err != nil {
		return nil, 0, err
	}
	return issues, total, nil
}

func (r *logRepository) GetIssueByFingerprint(fingerprint string, ownerUserID uint) (*domain.Issue, error) {
	var issue domain.Issue
	query := r.db.Model(&domain.Issue{}).Where("fingerprint = ?", fingerprint)
	if ownerUserID > 0 {
		query = query.Joins("JOIN sources ON sources.id = issues.source_id").Where("sources.user_id = ?", ownerUserID)
	}
	if err := query.First(&issue).Error; err != nil {
		return nil, err
	}
	return &issue, nil
}

func (r *logRepository) UpdateIssueStatus(fingerprint string, status string, ownerUserID uint) (*domain.Issue, error) {
	issue, err := r.GetIssueByFingerprint(fingerprint, ownerUserID)
	if err != nil {
		return nil, err
	}
	issue.Status = status
	if err := r.db.Save(issue).Error; err != nil {
		return nil, err
	}
	return issue, nil
}

func (r *logRepository) ListIssueLogs(limit int, offset int, fingerprint string, ownerUserID uint) ([]domain.LogEntry, int64, error) {
	var logs []domain.LogEntry
	var total int64
	query := r.db.Model(&domain.LogEntry{}).Where("fingerprint = ?", fingerprint)
	if ownerUserID > 0 {
		query = query.Joins("JOIN sources ON sources.id = log_entries.source_id").Where("sources.user_id = ?", ownerUserID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("log_entries.created_at desc").Limit(limit).Offset(offset).Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}
