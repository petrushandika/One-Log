package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/petrushandika/one-log/internal/service"
	"github.com/petrushandika/one-log/pkg/utils"
)

type ActivityMonitorHandler struct {
	service service.ActivityMonitorService
}

func NewActivityMonitorHandler(svc service.ActivityMonitorService) *ActivityMonitorHandler {
	return &ActivityMonitorHandler{service: svc}
}

// GET /api/activity/feed
func (h *ActivityMonitorHandler) GetActivityFeed(c *gin.Context) {
	sourceID := c.Query("source_id")
	userID := c.Query("user_id")
	action := c.Query("action")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var from, to time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		from, _ = time.Parse(time.RFC3339, fromStr)
	}
	if toStr := c.Query("to"); toStr != "" {
		to, _ = time.Parse(time.RFC3339, toStr)
	}

	feeds, total, err := h.service.GetActivityFeed(sourceID, userID, action, from, to, page, limit)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch activity feed", err.Error())
		return
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	utils.Success(c, http.StatusOK, "Activity feed retrieved successfully", gin.H{
		"items": feeds,
		"meta": gin.H{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
		},
	})
}

// GET /api/activity/top-users
func (h *ActivityMonitorHandler) GetTopActiveUsers(c *gin.Context) {
	sourceID := c.Query("source_id")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, err := h.service.GetTopActiveUsers(sourceID, days, limit)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch top active users", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Top active users retrieved successfully", users)
}

// GET /api/activity/by-resource
func (h *ActivityMonitorHandler) GetActivityByResource(c *gin.Context) {
	sourceID := c.Query("source_id")
	resourceType := c.Query("resource_type")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	activities, err := h.service.GetActivityByResource(sourceID, resourceType, days)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch activity by resource", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Activity by resource retrieved successfully", activities)
}

// GET /api/activity/users/:user_id/profile
func (h *ActivityMonitorHandler) GetUserProfile(c *gin.Context) {
	userID := c.Param("user_id")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	profile, err := h.service.GetUserProfile(userID, days)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch user profile", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "User profile retrieved successfully", profile)
}

// POST /api/activity/compliance-export
func (h *ActivityMonitorHandler) RequestComplianceExport(c *gin.Context) {
	var req struct {
		SourceID string    `json:"source_id" binding:"required"`
		Format   string    `json:"format" binding:"required,oneof=PDF CSV"`
		DateFrom time.Time `json:"date_from" binding:"required"`
		DateTo   time.Time `json:"date_to" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Get current user from context
	createdBy := "admin" // TODO: Get from JWT

	export, err := h.service.RequestComplianceExport(req.SourceID, req.Format, req.DateFrom, req.DateTo, createdBy)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to request compliance export", err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, "Compliance export requested successfully", export)
}

// GET /api/activity/compliance-exports
func (h *ActivityMonitorHandler) GetComplianceExports(c *gin.Context) {
	sourceID := c.Query("source_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	exports, total, err := h.service.GetComplianceExports(sourceID, page, limit)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch compliance exports", err.Error())
		return
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	utils.Success(c, http.StatusOK, "Compliance exports retrieved successfully", gin.H{
		"items": exports,
		"meta": gin.H{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
		},
	})
}
