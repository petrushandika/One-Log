package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/petrushandika/one-log/internal/repository"
	"github.com/petrushandika/one-log/pkg/utils"
)

type StatusHandler struct {
	sourceRepo repository.SourceRepository
}

func NewStatusHandler(repo repository.SourceRepository) *StatusHandler {
	return &StatusHandler{sourceRepo: repo}
}

// GET /api/v1/status (public)
func (h *StatusHandler) PublicStatus(c *gin.Context) {
	sources, err := h.sourceRepo.FindAll(0)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch status", err.Error())
		return
	}
	utils.Success(c, http.StatusOK, "Status retrieved successfully", gin.H{
		"sources": sources,
		"total":   len(sources),
	})
}
