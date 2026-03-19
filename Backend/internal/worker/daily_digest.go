package worker

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/petrushandika/one-log/internal/repository"
	"github.com/petrushandika/one-log/internal/service"
	"github.com/petrushandika/one-log/pkg/ai"
)

// DailyDigestWorker generates AI-powered daily summaries
type DailyDigestWorker struct {
	logRepo   repository.LogRepository
	aiClient  ai.GroqClient
	notifySvc service.NotificationService
}

func NewDailyDigestWorker(logRepo repository.LogRepository, aiClient ai.GroqClient, notifySvc service.NotificationService) *DailyDigestWorker {
	return &DailyDigestWorker{
		logRepo:   logRepo,
		aiClient:  aiClient,
		notifySvc: notifySvc,
	}
}

func (w *DailyDigestWorker) Start() {
	log.Println("Starting Daily Digest Worker...")
	go w.run()
}

func (w *DailyDigestWorker) run() {
	// Schedule for 8 AM every day
	for {
		now := time.Now()
		nextRun := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location())
		if now.After(nextRun) {
			nextRun = nextRun.Add(24 * time.Hour)
		}

		time.Sleep(time.Until(nextRun))
		w.generateDailyDigest()
	}
}

func (w *DailyDigestWorker) generateDailyDigest() {
	log.Println("Generating daily digest...")

	// Get yesterday's stats
	yesterday := time.Now().AddDate(0, 0, -1)
	startOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Get error stats
	errors, _, _ := w.logRepo.FindAll(100, 0, "", "ERROR", "", 0, &startOfDay, &endOfDay)
	criticals, _, _ := w.logRepo.FindAll(100, 0, "", "CRITICAL", "", 0, &startOfDay, &endOfDay)

	// Get top issues
	trend, _ := w.logRepo.GetErrorRateTrend(1, "", 0)

	// Build summary for AI
	summary := fmt.Sprintf("Daily Summary for %s:\n", yesterday.Format("2006-01-02"))
	summary += fmt.Sprintf("- Total Errors: %d\n", len(errors))
	summary += fmt.Sprintf("- Total Critical: %d\n", len(criticals))
	summary += fmt.Sprintf("- Error Trend: %v\n", trend)

	// Generate AI summary
	systemPrompt := `You are a system monitoring assistant. Generate daily digest reports based on log data. Be concise and actionable.`
	userPrompt := fmt.Sprintf(`Generate a daily digest report for a system monitoring platform.

Data:
%s

Please provide:
1. Executive Summary (2-3 sentences)
2. Key Issues (bullet points)
3. Recommendations (actionable items)
4. System Health Score (0-100)

Format in Markdown.`, summary)

	aiResponse, err := w.aiClient.AnalyzeLog(systemPrompt, userPrompt)
	if err != nil {
		log.Printf("Error generating AI digest: %v", err)
		return
	}

	// Send email
	subject := fmt.Sprintf("📊 Daily Digest - %s", yesterday.Format("2006-01-02"))
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail != "" {
		htmlBody := fmt.Sprintf(`
<h1>Daily System Digest</h1>
<p><strong>Date:</strong> %s</p>
<hr/>
%s
`, yesterday.Format("2006-01-02"), markdownToHTML(aiResponse))

		w.notifySvc.SendEmail(adminEmail, subject, htmlBody)
		log.Println("Daily digest email sent")
	}
}

func markdownToHTML(md string) string {
	// Simple markdown to HTML conversion
	// In production, use a proper markdown library
	return md
}
