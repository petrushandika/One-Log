package worker

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/internal/repository"
	"github.com/petrushandika/one-log/internal/service"
	"github.com/petrushandika/one-log/pkg/email"
)

type UptimeWorker struct {
	sourceRepo   repository.SourceRepository
	incidentRepo repository.IncidentRepository
	logSvc       service.LogService
	emailClient  *email.SMTPEmailService
}

func NewUptimeWorker(sourceRepo repository.SourceRepository, incidentRepo repository.IncidentRepository, logSvc service.LogService) *UptimeWorker {
	return &UptimeWorker{
		sourceRepo:   sourceRepo,
		incidentRepo: incidentRepo,
		logSvc:       logSvc,
		emailClient:  email.NewSMTPEmailService(),
	}
}

func (w *UptimeWorker) Start() {
	fmt.Println("[UptimeWorker] Started Uptime Monitoring Worker")
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			w.runCheck()
		}
	}()
}

func (w *UptimeWorker) runCheck() {
	sources, err := w.sourceRepo.FindAll(0)
	if err != nil {
		fmt.Printf("[UptimeWorker] Error fetching sources: %v\n", err)
		return
	}

	for i := range sources {
		if sources[i].HealthURL == "" {
			continue
		}
		// Goroutine to avoid blocking due to HTTP timeouts
		go w.ping(&sources[i])
	}
}

func (w *UptimeWorker) ping(s *domain.Source) {
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(s.HealthURL)

	isDown := err != nil || resp.StatusCode >= 400
	newStatus := "ONLINE"
	if isDown {
		newStatus = "OFFLINE"
	}

	// Only update and ingest log if status changes to prevent spamming
	if s.Status != newStatus {
		oldStatus := s.Status
		s.Status = newStatus
		_ = w.sourceRepo.Update(s)

		level := "INFO"
		if isDown {
			level = "CRITICAL" // Will trigger automatic notification & AI Analysis
			// Phase 4: Auto-create incident
			w.createIncident(s)
		} else if oldStatus == "OFFLINE" && newStatus == "ONLINE" {
			// Phase 4: Resolve incident and send recovery email
			w.resolveIncident(s)
		}

		_ = w.logSvc.IngestLog(domain.IngestLogRequest{
			Category: "SYSTEM_ERROR",
			Level:    level,
			Message:  fmt.Sprintf("Source '%s' Health Status changed from %s to %s", s.Name, oldStatus, newStatus),
			Context: map[string]interface{}{
				"health_url": s.HealthURL,
				"status":     newStatus,
			},
		}, s.ID)

		fmt.Printf("[UptimeWorker] Source %s Status updated to %s\n", s.Name, newStatus)
	}
}

// createIncident creates a new incident record when source goes down
func (w *UptimeWorker) createIncident(s *domain.Source) {
	// Check if there's already an open incident for this source
	existing, err := w.incidentRepo.FindOpenBySource(s.ID)
	if err == nil && existing != nil {
		// Incident already exists, don't create duplicate
		return
	}

	incident := &domain.Incident{
		SourceID:  s.ID,
		Status:    "OPEN",
		StartedAt: time.Now(),
		Message:   fmt.Sprintf("Source '%s' is DOWN. Health check failed.", s.Name),
	}

	if err := w.incidentRepo.Create(incident); err != nil {
		fmt.Printf("[UptimeWorker] Failed to create incident for source %s: %v\n", s.Name, err)
	} else {
		fmt.Printf("[UptimeWorker] Incident created for source %s (ID: %d)\n", s.Name, incident.ID)
	}
}

// resolveIncident resolves the open incident and sends recovery email
func (w *UptimeWorker) resolveIncident(s *domain.Source) {
	// Find open incident for this source
	incident, err := w.incidentRepo.FindOpenBySource(s.ID)
	if err != nil {
		fmt.Printf("[UptimeWorker] No open incident found for source %s\n", s.Name)
		return
	}

	// Resolve the incident
	message := fmt.Sprintf("Source '%s' is back ONLINE.", s.Name)
	if err := w.incidentRepo.Resolve(incident.ID, message); err != nil {
		fmt.Printf("[UptimeWorker] Failed to resolve incident for source %s: %v\n", s.Name, err)
		return
	}

	// Calculate downtime duration
	duration := time.Since(incident.StartedAt)
	durationStr := formatDuration(duration)

	fmt.Printf("[UptimeWorker] Incident resolved for source %s. Downtime: %s\n", s.Name, durationStr)

	// Send recovery email
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail != "" {
		if err := w.emailClient.SendRecoveryEmail(adminEmail, s.Name, durationStr); err != nil {
			fmt.Printf("[UptimeWorker] Failed to send recovery email: %v\n", err)
		} else {
			fmt.Printf("[UptimeWorker] Recovery email sent to %s\n", adminEmail)
		}
	}
}

// formatDuration formats a duration into a human-readable string
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0f seconds", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.1f minutes", d.Minutes())
	}
	return fmt.Sprintf("%.1f hours", d.Hours())
}
