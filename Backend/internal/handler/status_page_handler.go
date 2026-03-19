package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/petrushandika/one-log/internal/service"
	"github.com/petrushandika/one-log/pkg/utils"
)

type StatusPageHandler struct {
	service service.StatusPageService
}

func NewStatusPageHandler(svc service.StatusPageService) *StatusPageHandler {
	return &StatusPageHandler{service: svc}
}

// Admin Endpoints (JWT required)

// POST /api/admin/status-pages
func (h *StatusPageHandler) Create(c *gin.Context) {
	var req struct {
		SourceID    string `json:"source_id" binding:"required"`
		Slug        string `json:"slug" binding:"required"`
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		LogoURL     string `json:"logo_url"`
		IsPublic    bool   `json:"is_public"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	config, err := h.service.CreateStatusPage(req.SourceID, req.Slug, req.Title, req.Description, req.LogoURL, req.IsPublic)
	if err != nil {
		if err == service.ErrSlugExists {
			utils.Error(c, http.StatusConflict, "Slug already exists", err.Error())
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Failed to create status page", err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, "Status page created successfully", config)
}

// GET /api/admin/status-pages
func (h *StatusPageHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	configs, total, err := h.service.ListStatusPages(page, limit)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to list status pages", err.Error())
		return
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	utils.Success(c, http.StatusOK, "Status pages retrieved successfully", gin.H{
		"items": configs,
		"meta": gin.H{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
		},
	})
}

// GET /api/admin/status-pages/:source_id
func (h *StatusPageHandler) Get(c *gin.Context) {
	sourceID := c.Param("source_id")

	config, err := h.service.GetStatusPage(sourceID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "Status page not found", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Status page retrieved successfully", config)
}

// PATCH /api/admin/status-pages/:source_id
func (h *StatusPageHandler) Update(c *gin.Context) {
	sourceID := c.Param("source_id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	if err := h.service.UpdateStatusPage(sourceID, updates); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update status page", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Status page updated successfully", nil)
}

// DELETE /api/admin/status-pages/:source_id
func (h *StatusPageHandler) Delete(c *gin.Context) {
	sourceID := c.Param("source_id")

	if err := h.service.DeleteStatusPage(sourceID); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete status page", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Status page deleted successfully", nil)
}

// GET /api/admin/status-pages/:source_id/uptime
func (h *StatusPageHandler) GetUptime(c *gin.Context) {
	sourceID := c.Param("source_id")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	stats, err := h.service.GetUptimeStats(sourceID, days)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to get uptime stats", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Uptime stats retrieved successfully", stats)
}

// POST /api/admin/status-pages/:source_id/embed
func (h *StatusPageHandler) CreateEmbed(c *gin.Context) {
	sourceID := c.Param("source_id")

	embed, err := h.service.CreateEmbedWidget(sourceID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create embed widget", err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, "Embed widget created successfully", embed)
}

// Public Endpoints (No auth required for public status pages)

// GET /status/:slug
func (h *StatusPageHandler) PublicStatusPage(c *gin.Context) {
	slug := c.Param("slug")

	config, err := h.service.GetStatusPageBySlug(slug)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "Status page not found", err.Error())
		return
	}

	if !config.IsPublic {
		utils.Error(c, http.StatusForbidden, "Status page is private", nil)
		return
	}

	// Get uptime stats
	stats, _ := h.service.GetUptimeStats(config.SourceID, 30)
	allStats, _ := h.service.GetAllUptimeStats(30)

	utils.Success(c, http.StatusOK, "Status page retrieved", gin.H{
		"config":    config,
		"stats":     stats,
		"all_stats": allStats,
	})
}

// GET /embed/:token
func (h *StatusPageHandler) EmbedWidget(c *gin.Context) {
	token := c.Param("token")

	embed, err := h.service.GetEmbedWidget(token)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "Embed widget not found", err.Error())
		return
	}

	// Get uptime stats for the source
	stats, _ := h.service.GetUptimeStats(embed.SourceID, 30)

	utils.Success(c, http.StatusOK, "Embed widget data", gin.H{
		"embed": embed,
		"stats": stats,
	})
}
