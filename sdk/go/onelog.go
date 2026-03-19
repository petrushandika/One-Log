// Package onelog provides an official Go SDK for the One-Log (ULAM) platform.
//
// Example usage:
//
//	client := onelog.NewClient(
//	    onelog.WithAPIKey("ulam_live_xxxxxxxx"),
//	    onelog.WithEndpoint("https://api.ulam.your-domain.com"),
//	)
//
//	err := client.LogInfo(onelog.SystemError, "Database connection established", nil)
//	if err != nil {
//	    log.Printf("Failed to send log: %v", err)
//	}
package onelog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	defaultEndpoint = "http://localhost:8080/api/ingest"
	defaultTimeout  = 5 * time.Second
)

// LogLevel represents the severity level of a log
type LogLevel string

const (
	Critical LogLevel = "CRITICAL"
	Error    LogLevel = "ERROR"
	Warn     LogLevel = "WARN"
	Info     LogLevel = "INFO"
	Debug    LogLevel = "DEBUG"
)

// Category represents the type of log
type Category string

const (
	SystemError  Category = "SYSTEM_ERROR"
	UserActivity Category = "USER_ACTIVITY"
	AuthEvent    Category = "AUTH_EVENT"
	Performance  Category = "PERFORMANCE"
	Security     Category = "SECURITY"
	AuditTrail   Category = "AUDIT_TRAIL"
)

// LogEntry represents a log entry to be sent to One-Log
type LogEntry struct {
	Category   string                 `json:"category"`
	Level      string                 `json:"level"`
	Message    string                 `json:"message"`
	StackTrace string                 `json:"stack_trace,omitempty"`
	Context    map[string]interface{} `json:"context,omitempty"`
	IPAddress  string                 `json:"ip_address,omitempty"`
}

// Client is the One-Log SDK client
type Client struct {
	apiKey     string
	endpoint   string
	sourceID   string
	httpClient *http.Client
}

// ClientOption configures the Client
type ClientOption func(*Client)

// NewClient creates a new One-Log client
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		apiKey:   os.Getenv("ULAM_API_KEY"),
		endpoint: os.Getenv("ULAM_ENDPOINT"),
		sourceID: os.Getenv("ULAM_SOURCE_ID"),
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.endpoint == "" {
		c.endpoint = defaultEndpoint
	}

	return c
}

// WithAPIKey sets the API key
func WithAPIKey(key string) ClientOption {
	return func(c *Client) {
		c.apiKey = key
	}
}

// WithEndpoint sets the API endpoint
func WithEndpoint(endpoint string) ClientOption {
	return func(c *Client) {
		c.endpoint = endpoint
	}
}

// WithSourceID sets the source ID
func WithSourceID(sourceID string) ClientOption {
	return func(c *Client) {
		c.sourceID = sourceID
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// Log sends a log entry to One-Log
func (c *Client) Log(category Category, level LogLevel, message string, context map[string]interface{}) error {
	entry := LogEntry{
		Category:  string(category),
		Level:     string(level),
		Message:   message,
		Context:   context,
		IPAddress: c.getIPAddress(),
	}

	return c.send(entry)
}

// LogWithError logs an error with its stack trace
func (c *Client) LogWithError(category Category, level LogLevel, message string, err error, context map[string]interface{}) error {
	stackTrace := ""
	if err != nil {
		stackTrace = err.Error()
	}

	entry := LogEntry{
		Category:   string(category),
		Level:      string(level),
		Message:    message,
		StackTrace: stackTrace,
		Context:    context,
		IPAddress:  c.getIPAddress(),
	}

	return c.send(entry)
}

// LogCritical logs a critical level message
func (c *Client) LogCritical(category Category, message string, context map[string]interface{}) error {
	return c.Log(category, Critical, message, context)
}

// LogError logs an error level message
func (c *Client) LogError(category Category, message string, context map[string]interface{}) error {
	return c.Log(category, Error, message, context)
}

// LogWarn logs a warning level message
func (c *Client) LogWarn(category Category, message string, context map[string]interface{}) error {
	return c.Log(category, Warn, message, context)
}

// LogInfo logs an info level message
func (c *Client) LogInfo(category Category, message string, context map[string]interface{}) error {
	return c.Log(category, Info, message, context)
}

// LogDebug logs a debug level message
func (c *Client) LogDebug(category Category, message string, context map[string]interface{}) error {
	return c.Log(category, Debug, message, context)
}

// LogPerformance logs a performance metric
func (c *Client) LogPerformance(endpoint string, durationMs int, method string, statusCode int) error {
	return c.Log(Performance, Info, "API request completed", map[string]interface{}{
		"endpoint":    endpoint,
		"duration_ms": durationMs,
		"method":      method,
		"status_code": statusCode,
	})
}

// LogAuthEvent logs an authentication event
func (c *Client) LogAuthEvent(eventType string, authMethod string, userID string, success bool) error {
	level := Info
	if !success {
		level = Warn
	}

	return c.Log(AuthEvent, level, "Authentication event", map[string]interface{}{
		"event_type":  eventType,
		"auth_method": authMethod,
		"user_id":     userID,
		"success":     success,
	})
}

// LogAudit logs an audit trail event
func (c *Client) LogAudit(action string, actorID string, resourceType string, resourceID string, before interface{}, after interface{}) error {
	return c.Log(AuditTrail, Info, "Audit event", map[string]interface{}{
		"action":        action,
		"actor_id":      actorID,
		"resource_type": resourceType,
		"resource_id":   resourceID,
		"before":        before,
		"after":         after,
	})
}

// send sends the log entry to the API
func (c *Client) send(entry LogEntry) error {
	if c.apiKey == "" {
		return fmt.Errorf("API key is required")
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	req, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send log: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// getIPAddress returns the client IP address
func (c *Client) getIPAddress() string {
	// In a real implementation, this would detect the actual IP
	// For now, return empty to let the server detect it
	return ""
}

// FireAndForget logs asynchronously without waiting for response
// Use this when you don't want to block the main execution
func (c *Client) FireAndForget(category Category, level LogLevel, message string, context map[string]interface{}) {
	go func() {
		if err := c.Log(category, level, message, context); err != nil {
			// Silently ignore errors in fire-and-forget mode
			// In production, you might want to log this to stderr
		}
	}()
}
