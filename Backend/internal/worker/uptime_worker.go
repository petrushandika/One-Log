package worker

import (
	"fmt"
	"net/http"
	"time"

	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/internal/repository"
	"github.com/petrushandika/one-log/internal/service"
)

type UptimeWorker struct {
	sourceRepo repository.SourceRepository
	logSvc     service.LogService
}

func NewUptimeWorker(sourceRepo repository.SourceRepository, logSvc service.LogService) *UptimeWorker {
	return &UptimeWorker{
		sourceRepo: sourceRepo,
		logSvc:     logSvc,
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
	sources, err := w.sourceRepo.FindAll()
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
