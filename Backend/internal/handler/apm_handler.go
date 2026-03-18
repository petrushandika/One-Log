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

// GET /api/v1/apm/endpoints?period=24h&source_id=...
func (h *APMHandler) EndpointStats(c *gin.Context) {
	period := c.DefaultQuery("period", "24h")
	sourceID := c.Query("source_id")
	ownerUserID := c.GetUint("user_id")

	out, err := h.service.EndpointStats(period, sourceID, ownerUserID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}
	utils.Success(c, http.StatusOK, "APM endpoint stats retrieved successfully", out)
}
