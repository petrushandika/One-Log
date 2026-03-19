package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/petrushandika/one-log/internal/domain"
)

// TelegramService handles Telegram Bot notifications
type TelegramService struct {
	BotToken string
	ChatID   string
	Enabled  bool
}

// NewTelegramService creates a new Telegram service from environment variables
func NewTelegramService() *TelegramService {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	return &TelegramService{
		BotToken: botToken,
		ChatID:   chatID,
		Enabled:  botToken != "" && chatID != "",
	}
}

// SendAlert sends an alert message to Telegram for ERROR/CRITICAL logs
func (s *TelegramService) SendAlert(logEntry *domain.LogEntry) error {
	if !s.Enabled {
		return nil
	}

	// Build emoji based on level
	emoji := "⚠️"
	if logEntry.Level == "CRITICAL" {
		emoji = "🚨"
	}

	message := fmt.Sprintf(
		"%s *ULAM Alert: %s on %s*\n\n"+
			"*Category:* %s\n"+
			"*Level:* %s\n"+
			"*Message:* %s\n"+
			"*Time:* %s\n",
		emoji,
		logEntry.Level,
		logEntry.SourceID,
		logEntry.Category,
		logEntry.Level,
		escapeMarkdown(logEntry.Message),
		logEntry.CreatedAt.Format(time.RFC3339),
	)

	if logEntry.IPAddress != "" {
		message += fmt.Sprintf("*IP:* %s\n", logEntry.IPAddress)
	}

	if logEntry.StackTrace != "" {
		// Truncate stack trace for Telegram
		stackPreview := logEntry.StackTrace
		if len(stackPreview) > 500 {
			stackPreview = stackPreview[:500] + "..."
		}
		message += fmt.Sprintf("\n*Stack Trace:*\n```\n%s\n```", escapeMarkdown(stackPreview))
	}

	return s.sendMessage(message)
}

// SendRecoveryAlert sends a recovery notification when source comes back online
func (s *TelegramService) SendRecoveryAlert(sourceName string, downtimeDuration string) error {
	if !s.Enabled {
		return nil
	}

	message := fmt.Sprintf(
		"✅ *Server Recovered: %s*\n\n"+
			"Good news! The server has recovered and is back online.\n\n"+
			"*Downtime Duration:* %s\n"+
			"*Recovered At:* %s\n",
		sourceName,
		downtimeDuration,
		time.Now().Format(time.RFC3339),
	)

	return s.sendMessage(message)
}

// SendDailyDigest sends a daily summary of system status
func (s *TelegramService) SendDailyDigest(stats map[string]interface{}) error {
	if !s.Enabled {
		return nil
	}

	totalLogs := stats["total_logs"].(int64)
	errorCount := stats["error_count"].(int64)
	criticalCount := stats["critical_count"].(int64)

	message := fmt.Sprintf(
		"📊 *Daily System Report*\n\n"+
			"*Total Logs:* %d\n"+
			"*Errors:* %d\n"+
			"*Critical:* %d\n"+
			"*Report Date:* %s\n",
		totalLogs,
		errorCount,
		criticalCount,
		time.Now().Format("2006-01-02"),
	)

	return s.sendMessage(message)
}

// sendMessage sends a message to Telegram using Bot API
func (s *TelegramService) sendMessage(text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.BotToken)

	payload := map[string]interface{}{
		"chat_id":    s.ChatID,
		"text":       text,
		"parse_mode": "Markdown",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to send Telegram message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status %d", resp.StatusCode)
	}

	return nil
}

// escapeMarkdown escapes special characters for Telegram Markdown
func escapeMarkdown(text string) string {
	// Escape special markdown characters
	replacer := map[string]string{
		"_": "\\_",
		"*": "\\*",
		"[": "\\[",
		"]": "\\]",
		"(": "\\(",
		")": "\\)",
		"~": "\\~",
		"`": "\\`",
		">": "\\>",
		"#": "\\#",
		"+": "\\+",
		"-": "\\-",
		"=": "\\=",
		"|": "\\|",
		"{": "\\{",
		"}": "\\}",
		".": "\\.",
		"!": "\\!",
	}

	result := text
	for old, new := range replacer {
		result = replaceAll(result, old, new)
	}
	return result
}

// replaceAll replaces all occurrences without using strings package
func replaceAll(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result += new
			i += len(old)
		} else {
			result += string(s[i])
			i++
		}
	}
	return result
}
