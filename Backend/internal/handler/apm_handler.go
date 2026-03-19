package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/petrushandika/one-log/internal/service"
	"github.com/petrushandika/one-log/pkg/utils"
)

type APMHandler struct {
	service service.APMService
}

func NewAPMHandler(svc service.APMService) *APMHandler {
	return &APMHandler{service: svc}
}

// GET /api/apm/endpoints?period=24h&source_id=...
func (h *APMHandler) EndpointStats(c *gin.Context) {
	period := c.DefaultQuery("period", "24h")
	sourceID := c.Query("source_id")
	ownerUserID := c.GetUint("user_id")

	out, err := h.service.EndpointStats(period, sourceID, ownerUserID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid period", err.Error())
		return
	}
	utils.Success(c, http.StatusOK, "APM endpoint stats retrieved", out)
}

// GET /api/apm/timeline?period=24h&interval=1h&source_id=...&endpoint=...
func (h *APMHandler) ResponseTimeTimeline(c *gin.Context) {
	period := c.DefaultQuery("period", "24h")
	interval := c.DefaultQuery("interval", "1h")
	sourceID := c.Query("source_id")
	endpoint := c.Query("endpoint")
	ownerUserID := c.GetUint("user_id")

	out, err := h.service.ResponseTimeTimeline(period, interval, sourceID, endpoint, ownerUserID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Failed to get timeline", err.Error())
		return
	}
	utils.Success(c, http.StatusOK, "Response time timeline retrieved", out)
}
