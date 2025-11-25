package handlers

import (
	"job-scraper/internal/services/exporter"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func (h *handler) exportHandler(c *gin.Context) {
	exporter := exporter.NewExcelExporter(h.db)

	filename, err := exporter.ExportToExcel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	const contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	c.File(filename)

	os.Remove(filename)
}
