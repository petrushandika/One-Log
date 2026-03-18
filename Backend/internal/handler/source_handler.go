package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/internal/service"
	"github.com/petrushandika/one-log/pkg/utils"
)

type SourceHandler struct {
	service service.SourceService
}

func NewSourceHandler(service service.SourceService) *SourceHandler {
	return &SourceHandler{service: service}
}

func (h *SourceHandler) Create(c *gin.Context) {
	var req domain.CreateSourceRequest
	if err := h.shouldBindJSON(c, &req); err != nil {
		return
	}

	source, rawAPIKey, err := h.service.CreateSource(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create source", err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, "Source created successfully", gin.H{
		"id":      source.ID,
		"name":    source.Name,
		"api_key": rawAPIKey, // Only shown once upon creation!
	})
}

func (h *SourceHandler) GetAll(c *gin.Context) {
	sources, err := h.service.GetSources()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch sources", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Sources retrieved successfully", sources)
}

func (h *SourceHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	source, err := h.service.GetSourceByID(id)

	if err != nil {
		if err.Error() == "source not found" {
			utils.Error(c, http.StatusNotFound, "Source not found", nil)
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch source", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Source retrieved successfully", source)
}

func (h *SourceHandler) shouldBindJSON(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindJSON(req); err != nil {
		utils.Error(c, http.StatusUnprocessableEntity, "Validation failed", []utils.ErrorDetail{
			{Field: "body", Message: err.Error()},
		})
		return err
	}
	return nil
}

func (h *SourceHandler) RotateKey(c *gin.Context) {
	id := c.Param("id")
	rawAPIKey, err := h.service.RotateAPIKey(id)

	if err != nil {
		if err.Error() == "source not found" {
			utils.Error(c, http.StatusNotFound, "Source not found", nil)
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Failed to rotate API key", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "API Key rotated successfully", gin.H{
		"new_api_key": rawAPIKey, // Only shown once!
	})
}
