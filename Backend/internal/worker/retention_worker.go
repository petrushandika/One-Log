package worker

import (
	"log"
	"time"

	"github.com/petrushandika/one-log/internal/repository"
)

type RetentionWorker struct {
	repo repository.LogRepository
	days int
}

func NewRetentionWorker(repo repository.LogRepository, days int) *RetentionWorker {
	if days <= 0 {
		days = 30 // Default 30 days retention
	}
	return &RetentionWorker{
		repo: repo,
		days: days,
	}
}

// Start daily worker tick
func (w *RetentionWorker) Start() {
	log.Printf("[Worker] Starting Log Retention Worker (retention: %d days)", w.days)

	// Run immediately on start
	w.runCleanup()

	// Setup Ticker to run once every day
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for range ticker.C {
			w.runCleanup()
		}
	}()
}

func (w *RetentionWorker) runCleanup() {
	log.Println("[Worker] Running log cleanup sequence...")
	err := w.repo.DeleteOlderThan(w.days)
	if err != nil {
		log.Printf("[Worker] Failed cleanup: %v", err)
	} else {
		log.Println("[Worker] Cleanup successful")
	}
}
