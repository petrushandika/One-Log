package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/internal/handler"
	"gorm.io/datatypes"
)

// Mock Log Service for testing
type mockLogService struct {
	IngestLogCalled bool
	LastReq         domain.IngestLogRequest
	LastSourceID    string
}

func (m *mockLogService) IngestLog(req domain.IngestLogRequest, sourceID string) error {
	m.IngestLogCalled = true
	m.LastReq = req
	m.LastSourceID = sourceID
	return nil
}

func (m *mockLogService) GetLogs(limit int, page int, sourceID string, level string, category string, userID uint) ([]domain.LogEntry, int64, error) {
	return nil, 0, nil
}

func (m *mockLogService) GetLogByID(id uint) (*domain.LogEntry, error) {
	return &domain.LogEntry{
		ID:       id,
		SourceID: "some-uuid",
		Category: "SYSTEM_ERROR",
		Level:    "ERROR",
		Context:  datatypes.JSON(`{"key":"value"}`),
	}, nil
}

func (m *mockLogService) ManualAnalyzeLog(id uint) (*domain.LogEntry, error) {
	return &domain.LogEntry{
		ID:        id,
		Level:     "CRITICAL",
		AIInsight: datatypes.JSON(`{"analysis":"mock insight"}`),
	}, nil
}

func (m *mockLogService) GetStatsOverview(userID uint) (map[string]interface{}, error) {
	return map[string]interface{}{"total": 10}, nil
}

func (m *mockLogService) CheckBruteForce(ip string) (bool, error) {
	return false, nil
}

func TestIngestLog_Success(t *testing.T) {
	// Setup Ginkgo mode appropriately
	gin.SetMode(gin.TestMode)
	mService := &mockLogService{}
	logHandler := handler.NewLogHandler(mService)

	router := gin.Default()

	// Mock the middleware by directly setting source_id into context
	router.Use(func(c *gin.Context) {
		c.Set("source_id", "test-source-uuid-1234")
		c.Next()
	})

	router.POST("/api/v1/ingest", logHandler.Ingest)

	// Valid Payload
	payload := domain.IngestLogRequest{
		SourceID:   "test-source-uuid-1234",
		Category:   "SYSTEM_ERROR",
		Level:      "ERROR",
		Message:    "Database connection timed out",
		StackTrace: "at db.go:12\nat main.go:1",
		Context:    map[string]interface{}{"retry_count": 5},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/ingest", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Assert Status
	if recorder.Code != http.StatusAccepted {
		t.Errorf("Expected status code 202, got %v", recorder.Code)
	}

	// Assert Payload Parsed inside Service
	if !mService.IngestLogCalled {
		t.Errorf("Expected LogService.IngestLog to be called")
	}

	if mService.LastSourceID != "test-source-uuid-1234" {
		t.Errorf("Expected parsed source_id to match, got: %s", mService.LastSourceID)
	}

	if mService.LastReq.Message != "Database connection timed out" {
		t.Errorf("Expected message to be correctly passed, got: %v", mService.LastReq)
	}
}

func TestIngestLog_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mService := &mockLogService{}
	logHandler := handler.NewLogHandler(mService)

	router := gin.Default()
	router.POST("/api/v1/ingest", logHandler.Ingest)

	// Invalid Payload (Missing Category, Invalid Level)
	payload := map[string]interface{}{
		"level":   "BAD_LEVEL",
		"message": "",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/ingest", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Assert Status 422 Unprocessable Entity
	if recorder.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code 422, got %v", recorder.Code)
	}
}
