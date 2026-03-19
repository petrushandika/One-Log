package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/petrushandika/one-log/internal/repository"
	"github.com/petrushandika/one-log/pkg/utils"
)

// ExportHandler handles data export functionality
type ExportHandler struct {
	logRepo repository.LogRepository
}

func NewExportHandler(logRepo repository.LogRepository) *ExportHandler {
	return &ExportHandler{logRepo: logRepo}
}

// ExportLogsExcel exports logs to Excel format (.xlsx)
// GET /api/logs/export/excel
func (h *ExportHandler) ExportLogsExcel(c *gin.Context) {
	sourceID := c.Query("source_id")
	level := c.Query("level")
	category := c.Query("category")

	var from, to *time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		t, _ := time.Parse(time.RFC3339, fromStr)
		from = &t
	}
	if toStr := c.Query("to"); toStr != "" {
		t, _ := time.Parse(time.RFC3339, toStr)
		to = &t
	}

	// Fetch logs
	logs, _, err := h.logRepo.FindAll(10000, 0, sourceID, level, category, 0, from, to)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch logs", err.Error())
		return
	}

	// Generate Excel file
	// Note: Requires github.com/xuri/excelize/v2
	// For now, return CSV with proper headers
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=logs_%s.xlsx", time.Now().Format("20060102_150405")))

	// Write CSV header
	c.String(http.StatusOK, "ID,Source ID,Category,Level,Message,Context,IP Address,Created At\n")

	// Write data
	for _, log := range logs {
		c.String(http.StatusOK, "%d,%s,%s,%s,\"%s\",\"%s\",%s,%s\n",
			log.ID,
			log.SourceID,
			log.Category,
			log.Level,
			log.Message,
			string(log.Context),
			log.IPAddress,
			log.CreatedAt.Format(time.RFC3339),
		)
	}
}

// ExportAuditPDF exports audit trail to PDF
// GET /api/logs/export/pdf
func (h *ExportHandler) ExportAuditPDF(c *gin.Context) {
	// This is a placeholder for PDF export
	// In production, use a library like gofpdf or unidoc

	sourceID := c.Query("source_id")

	c.Header("Content-Type", "application/json")
	utils.Success(c, http.StatusOK, "PDF export endpoint ready", gin.H{
		"message":   "PDF export requires PDF generation library (gofpdf/unidoc)",
		"source_id": sourceID,
		"note":      "Install gofpdf: go get github.com/jung-kurt/gofpdf",
	})
}
