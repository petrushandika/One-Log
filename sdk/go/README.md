# One-Log Go SDK

Official Go SDK for One-Log (ULAM) - Unified Log & Activity Monitor

## Installation

```bash
go get github.com/petrushandika/one-log/sdk/go
```

## Quick Start

```go
package main

import (
    "log"
    "github.com/petrushandika/one-log/sdk/go/onelog"
)

func main() {
    // Create client
    client := onelog.NewClient(
        onelog.WithAPIKey("ulam_live_xxxxxxxx"),
        onelog.WithEndpoint("https://api.ulam.your-domain.com/api/ingest"),
    )

    // Send a log
    err := client.LogInfo(onelog.SystemError, "Application started", map[string]interface{}{
        "version": "1.0.0",
        "environment": "production",
    })
    if err != nil {
        log.Printf("Failed to send log: %v", err)
    }
}
```

## Configuration

### Environment Variables

```bash
export ULAM_API_KEY="ulam_live_xxxxxxxx"
export ULAM_ENDPOINT="https://api.ulam.your-domain.com/api/ingest"
export ULAM_SOURCE_ID="my-app"
```

### Client Options

```go
client := onelog.NewClient(
    onelog.WithAPIKey("your-api-key"),           // Required
    onelog.WithEndpoint("https://api..."),       // Optional (default: localhost)
    onelog.WithSourceID("my-app"),               // Optional
    onelog.WithHTTPClient(customHTTPClient),     // Optional
)
```

## Log Levels

```go
onelog.Critical  // System down, data loss
onelog.Error     // Errors that don't stop the app
onelog.Warn      // Warnings
onelog.Info      // General information
onelog.Debug     // Debug information
```

## Categories

```go
onelog.SystemError   // System errors, crashes
onelog.UserActivity  // User actions
onelog.AuthEvent     // Login/logout events
onelog.Performance   // Performance metrics
onelog.Security      // Security events
onelog.AuditTrail    // Audit logs
```

## Usage Examples

### Basic Logging

```go
// Simple info log
client.LogInfo(onelog.SystemError, "Database connection established", nil)

// Error log
client.LogError(onelog.SystemError, "Failed to connect to database", map[string]interface{}{
    "host": "localhost",
    "port": 5432,
})

// Critical log
client.LogCritical(onelog.SystemError, "System is out of memory", map[string]interface{}{
    "memory_usage": "95%",
})
```

### Logging with Error Details

```go
err := someOperation()
if err != nil {
    client.LogWithError(
        onelog.SystemError,
        onelog.Error,
        "Operation failed",
        err,
        map[string]interface{}{
            "operation": "data_sync",
            "user_id": "12345",
        },
    )
}
```

### Performance Logging

```go
start := time.Now()
// ... your operation ...
duration := time.Since(start)

client.LogPerformance("/api/users", int(duration.Milliseconds()), "GET", 200)
```

### Authentication Events

```go
// Successful login
client.LogAuthEvent("login_success", "google_oauth", "user_123", true)

// Failed login
client.LogAuthEvent("login_failed", "system_password", "user_123", false)
```

### Audit Trail

```go
client.LogAudit(
    "update",                                    // action
    "admin_001",                                 // actor
    "user",                                      // resource type
    "user_123",                                  // resource ID
    map[string]string{"name": "Old Name"},       // before
    map[string]string{"name": "New Name"},       // after
)
```

### Fire-and-Forget (Non-blocking)

Use when you don't want to block execution:

```go
// This runs asynchronously
client.FireAndForget(onelog.Performance, onelog.Info, "Request completed", map[string]interface{}{
    "endpoint": "/api/data",
    "duration_ms": 150,
})
```

## Complete Example

```go
package main

import (
    "log"
    "time"
    "github.com/petrushandika/one-log/sdk/go/onelog"
)

func main() {
    client := onelog.NewClient()
    
    // Log application start
    client.LogInfo(onelog.SystemError, "Application started", map[string]interface{}{
        "version": "1.0.0",
    })
    
    // Simulate work
    start := time.Now()
    
    // Some operation
    err := performOperation()
    if err != nil {
        client.LogError(onelog.SystemError, "Operation failed", map[string]interface{}{
            "error": err.Error(),
        })
    }
    
    // Log performance
    duration := time.Since(start)
    client.LogPerformance("/api/operation", int(duration.Milliseconds()), "POST", 200)
    
    // Log user activity
    client.LogInfo(onelog.UserActivity, "User performed action", map[string]interface{}{
        "user_id": "12345",
        "action": "export_data",
    })
}

func performOperation() error {
    // Your logic here
    return nil
}
```

## Error Handling

The SDK returns errors for:
- Missing API key
- Network failures
- API errors (non-202 status codes)

Always handle errors appropriately:

```go
err := client.LogInfo(onelog.SystemError, "Message", nil)
if err != nil {
    // Log to fallback (e.g., local file)
    log.Printf("Failed to send to One-Log: %v", err)
}
```

## Best Practices

1. **Use Fire-and-Forget for non-critical logs**: Don't block your app for logs
2. **Add context**: Include relevant IDs, timestamps, and metadata
3. **Use appropriate levels**: Don't log everything as ERROR
4. **Handle errors**: Always check for errors in critical paths
5. **Reuse client**: Create one client and reuse it (thread-safe)

## API Reference

### Client Methods

```go
// Basic logging
func (c *Client) Log(category Category, level LogLevel, message string, context map[string]interface{}) error

// Convenience methods
func (c *Client) LogCritical(category Category, message string, context map[string]interface{}) error
func (c *Client) LogError(category Category, message string, context map[string]interface{}) error
func (c *Client) LogWarn(category Category, message string, context map[string]interface{}) error
func (c *Client) LogInfo(category Category, message string, context map[string]interface{}) error
func (c *Client) LogDebug(category Category, message string, context map[string]interface{}) error

// Specialized logging
func (c *Client) LogWithError(category Category, level LogLevel, message string, err error, context map[string]interface{}) error
func (c *Client) LogPerformance(endpoint string, durationMs int, method string, statusCode int) error
func (c *Client) LogAuthEvent(eventType string, authMethod string, userID string, success bool) error
func (c *Client) LogAudit(action string, actorID string, resourceType string, resourceID string, before interface{}, after interface{}) error

// Async logging
func (c *Client) FireAndForget(category Category, level LogLevel, message string, context map[string]interface{})
```

## License

MIT License - See LICENSE file for details
