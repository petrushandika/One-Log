package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/petrushandika/one-log/internal/service"
	"github.com/petrushandika/one-log/pkg/utils"
)

type ActivityHandler struct {
	service service.ActivityService
}

func NewActivityHandler(svc service.ActivityService) *ActivityHandler {
	return &ActivityHandler{service: svc}
}

// GET /api/v1/activity
func (h *ActivityHandler) List(c *gin.Context) {
	limit := c.DefaultQuery("limit", "20")
	page := c.DefaultQuery("page", "1")
	sourceID := c.Query("source_id")
	category := c.DefaultQuery("category", "AUTH_EVENT,USER_ACTIVITY,AUDIT_TRAIL")
	eventType := c.Query("event_type")
	authMethod := c.Query("auth_method")
	subjectUserID := c.Query("user_id")
	from := c.Query("from")
	to := c.Query("to")

	ownerUserID := c.GetUint("user_id")

	categories := []string{}
	for _, p := range strings.Split(category, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			categories = append(categories, p)
		}
	}

	items, meta, err := h.service.List(limit, page, sourceID, categories, eventType, authMethod, subjectUserID, from, to, ownerUserID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch activity", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Activity retrieved successfully", gin.H{
		"items": items,
		"meta":  meta,
	})
}

// GET /api/v1/activity/summary
func (h *ActivityHandler) Summary(c *gin.Context) {
	period := c.DefaultQuery("period", "7d")
	sourceID := c.Query("source_id")
	ownerUserID := c.GetUint("user_id")

	out, err := h.service.Summary(period, sourceID, ownerUserID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid period", err.Error())
		return
	}
	utils.Success(c, http.StatusOK, "Activity summary retrieved successfully", out)
}

// GET /api/v1/activity/users/:user_id
func (h *ActivityHandler) ByUser(c *gin.Context) {
	subjectUserID := c.Param("user_id")
	period := c.DefaultQuery("period", "30d")
	category := c.DefaultQuery("category", "AUTH_EVENT,USER_ACTIVITY,AUDIT_TRAIL")
	ownerUserID := c.GetUint("user_id")

	categories := []string{}
	for _, p := range strings.Split(category, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			categories = append(categories, p)
		}
	}

	out, err := h.service.ByUser(subjectUserID, period, categories, ownerUserID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}
	utils.Success(c, http.StatusOK, "User activity retrieved successfully", out)
}

// GET /api/v1/activity/suspicious
func (h *ActivityHandler) Suspicious(c *gin.Context) {
	period := c.DefaultQuery("period", "24h")
	sourceID := c.Query("source_id")
	limit := c.DefaultQuery("limit", "20")
	page := c.DefaultQuery("page", "1")
	ownerUserID := c.GetUint("user_id")

	items, meta, err := h.service.Suspicious(limit, page, period, sourceID, ownerUserID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}
	utils.Success(c, http.StatusOK, "Suspicious activity retrieved successfully", gin.H{
		"items": items,
		"meta":  meta,
	})
}

func parseTimeRFC3339(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
