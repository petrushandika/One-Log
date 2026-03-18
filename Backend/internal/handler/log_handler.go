package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/internal/service"
	"github.com/petrushandika/one-log/pkg/utils"
)

type LogHandler struct {
	service service.LogService
}

func NewLogHandler(service service.LogService) *LogHandler {
	return &LogHandler{service: service}
}

func (h *LogHandler) Ingest(c *gin.Context) {
	var req domain.IngestLogRequest

	// Try process JSON
	if err := h.shouldBindJSON(c, &req); err != nil {
		return
	}

	// Get SourceID from middleware
	sourceID, exists := c.Get("source_id")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "Invalid source session", nil)
		return
	}

	// Call Service
	if err := h.service.IngestLog(req, sourceID.(string)); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to ingest log", err.Error())
		return
	}

	utils.Success(c, http.StatusAccepted, "Log ingested successfully", nil)
}

func (h *LogHandler) shouldBindJSON(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindJSON(req); err != nil {
		utils.Error(c, http.StatusUnprocessableEntity, "Validation failed", []utils.ErrorDetail{
			{Field: "body", Message: err.Error()},
		})
		return err
	}
	return nil
}

// GetAll handles GET /api/v1/logs requests
func (h *LogHandler) GetAll(c *gin.Context) {
	// Simple query params parsing
	limitStr := c.DefaultQuery("limit", "20")
	pageStr := c.DefaultQuery("page", "1")
	sourceID := c.Query("source_id")
	level := c.Query("level")
	category := c.Query("category")

	var limit, page int
	// Parse basic using Sscanf, ignoring error and falling back to default inside service
	_, _ = fmt.Sscanf(limitStr, "%d", &limit)
	_, _ = fmt.Sscanf(pageStr, "%d", &page)

	logs, total, err := h.service.GetLogs(limit, page, sourceID, level, category)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch logs", err.Error())
		return
	}

	// Meta info for client pagination
	utils.Success(c, http.StatusOK, "Logs retrieved successfully", gin.H{
		"items": logs,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// GetByID handles GET /api/v1/logs/:id requests
func (h *LogHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid log ID format", nil)
		return
	}

	logEntry, err := h.service.GetLogByID(id)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "Log not found or an error occurred", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Log retrieved successfully", logEntry)
}

// Analyze handles POST /api/v1/logs/:id/analyze requests
func (h *LogHandler) Analyze(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid log ID format", nil)
		return
	}

	logEntry, err := h.service.ManualAnalyzeLog(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to analyze log with AI", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "AI Analysis completed", logEntry)
}

// GetStatsOverview handles GET /api/v1/stats/overview requests
func (h *LogHandler) GetStatsOverview(c *gin.Context) {
	stats, err := h.service.GetStatsOverview()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to load stats overview", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Stats retrieved successfully", stats)
}
