package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/petrushandika/one-log/internal/service"
	"github.com/petrushandika/one-log/pkg/utils"
)

type ConfigHandler struct {
	service service.ConfigService
}

func NewConfigHandler(s service.ConfigService) *ConfigHandler {
	return &ConfigHandler{service: s}
}

type saveConfigRequest struct {
	Key      string `json:"key" binding:"required"`
	Value    string `json:"value" binding:"required"`
	IsSecret bool   `json:"is_secret"`
}

func (h *ConfigHandler) Save(c *gin.Context) {
	sourceID := c.Param("id")
	var req saveConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.service.SaveConfig(sourceID, req.Key, req.Value, req.IsSecret); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to save config", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Config saved successfully", nil)
}

func (h *ConfigHandler) GetBySource(c *gin.Context) {
	sourceID := c.Param("id")
	configs, err := h.service.GetConfigsBySource(sourceID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to retrieve configs", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Configs retrieved successfully", configs)
}
