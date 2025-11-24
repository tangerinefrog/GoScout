package handlers

import (
	"job-scraper/internal/services/scraper"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ScrapeRequest struct {
	Keywords   string `json:"keywords"`
	PeriodDays int    `json:"period_days"`
}

func (h *Handler) scrapeHandler(c *gin.Context) {
	var req ScrapeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	keywords := req.Keywords
	if keywords == "" {
		c.String(http.StatusBadRequest, "Missing 'keywords' query param")
		return
	}

	periodDays := req.PeriodDays
	if periodDays == 0 {
		c.String(http.StatusBadRequest, "Missing 'period_days' query param")
		return
	}

	scraper := scraper.NewScraper(h.db)
	_, err := scraper.ScrapeLinkedInJobs(keywords, time.Duration(periodDays)*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
