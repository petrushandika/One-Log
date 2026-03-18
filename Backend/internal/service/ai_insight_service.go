package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/internal/repository"
	"github.com/petrushandika/one-log/pkg/ai"
	"gorm.io/datatypes"
)

type AIService interface {
	AnalyzeCriticalLog(logEntry *domain.LogEntry)
	ManualAnalyzeLog(id uint) (*domain.LogEntry, error)
}

type aiService struct {
	groqClient ai.GroqClient
	logRepo    repository.LogRepository
}

func NewAIService(repo repository.LogRepository) AIService {
	return &aiService{
		groqClient: ai.NewGroqClient(),
		logRepo:    repo,
	}
}

// AnalyzeCriticalLog is an automatic background task that fires for CRITICAL logs.
func (s *aiService) AnalyzeCriticalLog(logEntry *domain.LogEntry) {
	if logEntry.Level != "CRITICAL" {
		return
	}

	analysisResult, err := s.generateAIInsight(logEntry)
	if err != nil {
		log.Printf("[AI Service] Failed to analyze log %d: %v", logEntry.ID, err)
		return
	}

	// Update the Log in Database with AI Insight
	insightBytes, _ := json.Marshal(map[string]string{
		"analysis": analysisResult,
	})

	logEntry.AIInsight = datatypes.JSON(insightBytes)
	err = s.logRepo.Update(logEntry)
	if err != nil {
		log.Printf("[AI Service] Failed to save AI Insight to DB %d: %v", logEntry.ID, err)
	}
}

// ManualAnalyzeLog fires when a user clicks 'Analyze' from dashboard.
func (s *aiService) ManualAnalyzeLog(id uint) (*domain.LogEntry, error) {
	logEntry, err := s.logRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Only analyze if it hasn't been analyzed before
	if len(logEntry.AIInsight) == 0 || string(logEntry.AIInsight) == "null" {
		analysisResult, err := s.generateAIInsight(logEntry)
		if err != nil {
			return nil, err
		}

		insightBytes, _ := json.Marshal(map[string]string{
			"analysis": analysisResult,
		})
		logEntry.AIInsight = datatypes.JSON(insightBytes)

		err = s.logRepo.Update(logEntry)
		if err != nil {
			return nil, err
		}
	}

	return logEntry, nil
}

func (s *aiService) generateAIInsight(logEntry *domain.LogEntry) (string, error) {
	// Construct Context-Aware System Prompt
	systemPrompt := `You are an expert Security Analyst and System Forensics Engineer for ULAM (Unified Log & Activity Monitor).
Your job is to read log data and provide a concise, human-readable analysis.
Do not provide generic coding advice. Provide Root Cause Analysis and specific Mitigation steps.

CONTEXT AWARENESS:
If level is CRITICAL or ERROR, provide urgency and 3 steps fixing.
If level is INFO or WARN and Category is USER_ACTIVITY or AUTH_EVENT, this is mostly normal standard tracking. DO NOT suggest "Server shutdown" or panic. Describe what the user did objectively.
`

	userMessage := fmt.Sprintf(
		"Please analyze the following log:\nLevel: %s\nCategory: %s\nMessage: %s\nIP Address: %s\nContext: %s\nStack Trace: %s",
		logEntry.Level,
		logEntry.Category,
		logEntry.Message,
		logEntry.IPAddress,
		string(logEntry.Context),
		logEntry.StackTrace,
	)

	return s.groqClient.AnalyzeLog(systemPrompt, userMessage)
}
