package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/petrushandika/one-log/internal/service"
	"github.com/petrushandika/one-log/pkg/utils"
)

type IssueHandler struct {
	service service.IssueService
}

func NewIssueHandler(svc service.IssueService) *IssueHandler {
	return &IssueHandler{service: svc}
}

// GET /api/v1/issues
func (h *IssueHandler) List(c *gin.Context) {
	limit := c.DefaultQuery("limit", "20")
	page := c.DefaultQuery("page", "1")
	sourceID := c.Query("source_id")
	status := c.Query("status")
	ownerUserID := c.GetUint("user_id")

	items, meta, err := h.service.List(limit, page, sourceID, status, ownerUserID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch issues", err.Error())
		return
	}
	utils.Success(c, http.StatusOK, "Issues retrieved successfully", gin.H{
		"items": items,
		"meta":  meta,
	})
}

// GET /api/v1/issues/:fingerprint
func (h *IssueHandler) Get(c *gin.Context) {
	fp := c.Param("fingerprint")
	ownerUserID := c.GetUint("user_id")
	item, err := h.service.Get(fp, ownerUserID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "Issue not found", err.Error())
		return
	}
	utils.Success(c, http.StatusOK, "Issue retrieved successfully", item)
}

// PATCH /api/v1/issues/:fingerprint
func (h *IssueHandler) UpdateStatus(c *gin.Context) {
	fp := c.Param("fingerprint")
	ownerUserID := c.GetUint("user_id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusUnprocessableEntity, "Validation failed", []utils.ErrorDetail{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	item, err := h.service.UpdateStatus(fp, req.Status, ownerUserID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Failed to update issue status", err.Error())
		return
	}
	utils.Success(c, http.StatusOK, "Issue status updated successfully", item)
}

// GET /api/v1/issues/:fingerprint/logs
func (h *IssueHandler) Logs(c *gin.Context) {
	fp := c.Param("fingerprint")
	limit := c.DefaultQuery("limit", "20")
	page := c.DefaultQuery("page", "1")
	ownerUserID := c.GetUint("user_id")

	items, meta, err := h.service.Logs(fp, limit, page, ownerUserID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch issue logs", err.Error())
		return
	}
	utils.Success(c, http.StatusOK, "Issue logs retrieved successfully", gin.H{
		"items": items,
		"meta":  meta,
	})
}
