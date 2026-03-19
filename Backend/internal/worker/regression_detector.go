package worker

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/petrushandika/one-log/internal/repository"
	"github.com/petrushandika/one-log/internal/service"
)

// RegressionDetector checks for issues that were RESOLVED but reoccurred
type RegressionDetector struct {
	logRepo       repository.LogRepository
	notifySvc     service.NotificationService
	checkInterval time.Duration
}

func NewRegressionDetector(logRepo repository.LogRepository, notifySvc service.NotificationService) *RegressionDetector {
	return &RegressionDetector{
		logRepo:       logRepo,
		notifySvc:     notifySvc,
		checkInterval: 5 * time.Minute,
	}
}

func (w *RegressionDetector) Start() {
	log.Println("Starting Regression Detector Worker...")
	go w.run()
}

func (w *RegressionDetector) run() {
	ticker := time.NewTicker(w.checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		w.detectRegressions()
	}
}

func (w *RegressionDetector) detectRegressions() {
	// Find issues that:
	// 1. Were RESOLVED
	// 2. Have new occurrences (LastSeenAt > ResolvedAt)
	// 3. Haven't had regression alert sent yet

	issues, err := w.logRepo.FindResolvedIssuesWithNewOccurrences()
	if err != nil {
		log.Printf("Error finding regression issues: %v", err)
		return
	}

	for _, issue := range issues {
		if issue.RegressionAlertSent {
			continue
		}

		// Mark as regression
		w.logRepo.MarkAsRegression(issue.Fingerprint)

		// Send alert
		w.sendRegressionAlert(&issue)

		log.Printf("Regression detected: Issue %s has reoccurred", issue.Fingerprint)
	}
}

func (w *RegressionDetector) sendRegressionAlert(issue *repository.IssueWithSource) {
	subject := "🚨 Regression Alert: Resolved Issue Has Reoccurred"
	body := fmt.Sprintf(`
<h2>Regression Detected</h2>
<p>An issue that was previously marked as RESOLVED has reoccurred:</p>

<ul>
  <li><strong>Issue:</strong> %s</li>
  <li><strong>Source:</strong> %s</li>
  <li><strong>Level:</strong> %s</li>
  <li><strong>New Occurrences:</strong> %d</li>
  <li><strong>Last Seen:</strong> %s</li>
</ul>

<p>The issue status has been automatically changed to OPEN.</p>

<p><a href="http://localhost:5173/issues">View in Dashboard</a></p>
`, issue.MessageSample, issue.SourceName, issue.Level, issue.OccurrenceCount, issue.LastSeenAt.Format("2006-01-02 15:04:05"))

	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail != "" {
		w.notifySvc.SendEmail(adminEmail, subject, body)
	}
}
