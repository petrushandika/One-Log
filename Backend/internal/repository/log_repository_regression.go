package repository

import (
	"time"

	"github.com/petrushandika/one-log/internal/domain"
)

// FindResolvedIssuesWithNewOccurrences finds issues that were resolved but have new log entries
func (r *logRepository) FindResolvedIssuesWithNewOccurrences() ([]IssueWithSource, error) {
	var issues []IssueWithSource

	query := `
		SELECT 
			i.fingerprint,
			i.source_id,
			s.name as source_name,
			i.level,
			i.message_sample,
			i.occurrence_count,
			i.last_seen_at,
			i.resolved_at,
			i.regression_alert_sent
		FROM issues i
		JOIN sources s ON s.id = i.source_id
		WHERE i.status = 'RESOLVED'
			AND i.resolved_at IS NOT NULL
			AND i.last_seen_at > i.resolved_at
			AND i.regression_alert_sent = false
	`

	rows, err := r.db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var issue IssueWithSource
		var resolvedAt *time.Time

		if err := rows.Scan(
			&issue.Fingerprint,
			&issue.SourceID,
			&issue.SourceName,
			&issue.Level,
			&issue.MessageSample,
			&issue.OccurrenceCount,
			&issue.LastSeenAt,
			&resolvedAt,
			&issue.RegressionAlertSent,
		); err != nil {
			continue
		}
		issue.ResolvedAt = resolvedAt
		issues = append(issues, issue)
	}

	return issues, nil
}

// MarkAsRegression marks an issue as regression and changes status back to OPEN
func (r *logRepository) MarkAsRegression(fingerprint string) error {
	return r.db.Model(&domain.Issue{}).Where("fingerprint = ?", fingerprint).Updates(map[string]interface{}{
		"status":                "OPEN",
		"is_regression":         true,
		"regression_alert_sent": true,
		"updated_at":            time.Now(),
	}).Error
}
