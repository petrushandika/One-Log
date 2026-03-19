package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/petrushandika/one-log/internal/service"
	"github.com/petrushandika/one-log/pkg/utils"
)

type APMThresholdHandler struct {
	service service.APMThresholdService
}

func NewAPMThresholdHandler(svc service.APMThresholdService) *APMThresholdHandler {
	return &APMThresholdHandler{service: svc}
}

// GET /api/apm/thresholds
func (h *APMThresholdHandler) List(c *gin.Context) {
	sourceID := c.Query("source_id")

	thresholds, err := h.service.List(sourceID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch thresholds", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Thresholds retrieved successfully", thresholds)
}

// GET /api/apm/thresholds/:id
func (h *APMThresholdHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	threshold, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, "Threshold not found", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Threshold retrieved successfully", threshold)
}

// POST /api/apm/thresholds
func (h *APMThresholdHandler) Create(c *gin.Context) {
	var req struct {
		SourceID    string `json:"source_id" binding:"required"`
		Endpoint    string `json:"endpoint" binding:"required"`
		P95Limit    int    `json:"p95_limit"`
		P99Limit    int    `json:"p99_limit"`
		EmailNotify bool   `json:"email_notify"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	if req.P95Limit <= 0 {
		req.P95Limit = 1000
	}
	if req.P99Limit <= 0 {
		req.P99Limit = 2000
	}

	threshold, err := h.service.Create(req.SourceID, req.Endpoint, req.P95Limit, req.P99Limit, req.EmailNotify)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create threshold", err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, "Threshold created successfully", threshold)
}

// PATCH /api/apm/thresholds/:id
func (h *APMThresholdHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	var req struct {
		P95Limit    int  `json:"p95_limit"`
		P99Limit    int  `json:"p99_limit"`
		EmailNotify bool `json:"email_notify"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	if err := h.service.Update(uint(id), req.P95Limit, req.P99Limit, req.EmailNotify); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update threshold", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Threshold updated successfully", nil)
}

// DELETE /api/apm/thresholds/:id
func (h *APMThresholdHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete threshold", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Threshold deleted successfully", nil)
}

// GET /api/apm/slow-queries
func (h *APMThresholdHandler) GetSlowQueries(c *gin.Context) {
	sourceID := c.Query("source_id")
	thresholdStr := c.DefaultQuery("threshold", "2000")

	threshold, _ := strconv.Atoi(thresholdStr)
	if threshold <= 0 {
		threshold = 2000
	}

	queries, err := h.service.GetSlowQueries(sourceID, threshold)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch slow queries", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Slow queries retrieved successfully", queries)
}

// GET /api/apm/slow-queries/trend
func (h *APMThresholdHandler) GetSlowQueryTrend(c *gin.Context) {
	sourceID := c.Query("source_id")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	trend, err := h.service.GetSlowQueryTrend(sourceID, days)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch slow query trend", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Slow query trend retrieved successfully", trend)
}

// GET /api/apm/apdex
func (h *APMThresholdHandler) GetApdexScore(c *gin.Context) {
	sourceID := c.Query("source_id")
	endpoint := c.Query("endpoint")
	thresholdMs, _ := strconv.Atoi(c.DefaultQuery("threshold_ms", "1000"))

	score, err := h.service.CalculateApdexScore(sourceID, endpoint, thresholdMs)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to calculate Apdex score", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Apdex score calculated successfully", score)
}

// GET /api/apm/threshold-alerts
func (h *APMThresholdHandler) GetThresholdAlerts(c *gin.Context) {
	alerts, err := h.service.CheckThresholdAlerts()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to check threshold alerts", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Threshold alerts retrieved successfully", alerts)
}
