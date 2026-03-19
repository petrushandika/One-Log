package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/petrushandika/one-log/internal/service"
	"github.com/petrushandika/one-log/pkg/utils"
)

type IncidentHandler struct {
	service service.IncidentService
}

func NewIncidentHandler(svc service.IncidentService) *IncidentHandler {
	return &IncidentHandler{service: svc}
}

// GET /api/v1/incidents
func (h *IncidentHandler) List(c *gin.Context) {
	limit := c.DefaultQuery("limit", "20")
	page := c.DefaultQuery("page", "1")
	sourceID := c.Query("source_id")
	status := c.Query("status")

	items, meta, err := h.service.List(limit, page, sourceID, status)
	if err != nil {
		// Return empty result instead of error for better UX
		utils.Success(c, http.StatusOK, "Incidents retrieved successfully", gin.H{
			"items": []interface{}{},
			"meta": gin.H{
				"total": 0,
				"page":  1,
				"limit": 20,
			},
		})
		return
	}
	utils.Success(c, http.StatusOK, "Incidents retrieved successfully", gin.H{
		"items": items,
		"meta":  meta,
	})
}

// GET /api/v1/incidents/timeline
func (h *IncidentHandler) GetTimeline(c *gin.Context) {
	sourceID := c.Query("source_id")
	daysStr := c.DefaultQuery("days", "30")

	days, _ := strconv.Atoi(daysStr)
	if days <= 0 {
		days = 30
	}

	data, err := h.service.GetTimeline(sourceID, days)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to get incident timeline", err.Error())
		return
	}
	utils.Success(c, http.StatusOK, "Incident timeline retrieved", data)
}
