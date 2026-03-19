package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/petrushandika/one-log/internal/service"
	"github.com/petrushandika/one-log/pkg/utils"
)

type ActivityAnalyticsHandler struct {
	service service.ActivityAnalyticsService
}

func NewActivityAnalyticsHandler(svc service.ActivityAnalyticsService) *ActivityAnalyticsHandler {
	return &ActivityAnalyticsHandler{service: svc}
}

// GET /api/activity/analytics/methods
func (h *ActivityAnalyticsHandler) GetAuthMethodBreakdown(c *gin.Context) {
	sourceID := c.Query("source_id")
	daysStr := c.DefaultQuery("days", "30")

	days, _ := strconv.Atoi(daysStr)
	if days <= 0 {
		days = 30
	}

	data, err := h.service.GetAuthMethodBreakdown(sourceID, days)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to get auth method breakdown", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Auth method breakdown retrieved", data)
}

// GET /api/activity/analytics/timeline
func (h *ActivityAnalyticsHandler) GetLoginTimeline(c *gin.Context) {
	sourceID := c.Query("source_id")
	daysStr := c.DefaultQuery("days", "7")

	days, _ := strconv.Atoi(daysStr)
	if days <= 0 {
		days = 7
	}

	data, err := h.service.GetLoginTimeline(sourceID, days)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to get login timeline", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Login timeline retrieved", data)
}

// GET /api/activity/analytics/heatmap
func (h *ActivityAnalyticsHandler) GetFailedLoginHeatmap(c *gin.Context) {
	sourceID := c.Query("source_id")
	daysStr := c.DefaultQuery("days", "30")

	days, _ := strconv.Atoi(daysStr)
	if days <= 0 {
		days = 30
	}

	data, err := h.service.GetFailedLoginHeatmap(sourceID, days)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to get failed login heatmap", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Failed login heatmap retrieved", data)
}

// GET /api/activity/sessions
func (h *ActivityAnalyticsHandler) GetRecentSessions(c *gin.Context) {
	sourceID := c.Query("source_id")
	limit := c.DefaultQuery("limit", "20")
	page := c.DefaultQuery("page", "1")

	items, meta, err := h.service.GetRecentSessions(limit, page, sourceID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch sessions", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Sessions retrieved successfully", gin.H{
		"items": items,
		"meta":  meta,
	})
}
