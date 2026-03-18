package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/petrushandika/one-log/internal/service"
	"github.com/petrushandika/one-log/pkg/utils"
)

type ConfigHandler struct {
	service   service.ConfigService
	sourceSvc service.SourceService
}

func NewConfigHandler(s service.ConfigService, sourceSvc service.SourceService) *ConfigHandler {
	return &ConfigHandler{service: s, sourceSvc: sourceSvc}
}

type saveConfigRequest struct {
	Key         string `json:"key" binding:"required"`
	Value       string `json:"value" binding:"required"`
	IsSecret    bool   `json:"is_secret"`
	Environment string `json:"environment" binding:"omitempty"`
}

func (h *ConfigHandler) Save(c *gin.Context) {
	sourceID := c.Param("id")
	userID := c.GetUint("user_id")

	// Verify Ownership
	if _, err := h.sourceSvc.GetSourceByID(sourceID, userID); err != nil {
		utils.Error(c, http.StatusNotFound, "Source not found or access denied", "")
		return
	}

	var req saveConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.service.SaveConfig(sourceID, req.Environment, req.Key, req.Value, req.IsSecret, userID); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to save config", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Config saved successfully", nil)
}

func (h *ConfigHandler) GetBySource(c *gin.Context) {
	sourceID := c.Param("id")
	userID := c.GetUint("user_id")
	environment := c.Query("environment")
	reveal := c.Query("reveal") == "true"

	// Verify Ownership
	if _, err := h.sourceSvc.GetSourceByID(sourceID, userID); err != nil {
		utils.Error(c, http.StatusNotFound, "Source not found or access denied", "")
		return
	}

	configs, err := h.service.GetConfigsBySource(sourceID, environment, reveal)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to retrieve configs", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Configs retrieved successfully", configs)
}

// GET /api/v1/sources/:id/configs/history?environment=production&key=FOO&limit=50&reveal=true
func (h *ConfigHandler) History(c *gin.Context) {
	sourceID := c.Param("id")
	userID := c.GetUint("user_id")

	// Verify Ownership
	if _, err := h.sourceSvc.GetSourceByID(sourceID, userID); err != nil {
		utils.Error(c, http.StatusNotFound, "Source not found or access denied", "")
		return
	}

	environment := c.DefaultQuery("environment", "production")
	key := c.Query("key")
	limitStr := c.DefaultQuery("limit", "50")
	reveal := c.Query("reveal") == "true"

	var limit int
	_, _ = fmt.Sscanf(limitStr, "%d", &limit)

	history, err := h.service.GetHistory(sourceID, environment, key, limit, reveal)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to retrieve config history", err.Error())
		return
	}
	utils.Success(c, http.StatusOK, "Config history retrieved successfully", history)
}
