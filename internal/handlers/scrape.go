package handlers

import (
	"job-scraper/internal/services/scraper"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ScrapeRequest struct {
	SearchBy    string   `json:"search_by"`
	FilterBy    []string `json:"filter_by"`
	PeriodHours int      `json:"period_hours"`
}

func (h *handler) scrapeHandler(c *gin.Context) {
	var req ScrapeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	searchBy := req.SearchBy
	if searchBy == "" {
		c.String(http.StatusBadRequest, "Missing 'search_keywords' query param")
		return
	}

	periodDays := req.PeriodHours
	if periodDays == 0 {
		c.String(http.StatusBadRequest, "Missing 'period_hours' query param")
		return
	}

	scraper := scraper.NewScraper(h.db)
	_, err := scraper.ScrapeLinkedInJobs(c.Request.Context(), searchBy, req.FilterBy, time.Duration(periodDays)*time.Hour)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
