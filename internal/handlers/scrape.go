package handlers

import (
	"errors"
	"io"
	"job-scraper/internal/services/scraper"
	"log"
	"net/http"
	"strings"
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
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	searchBy := req.SearchBy
	if searchBy == "" {
		searchBy = h.config.SearchQuery
		if searchBy == "" {
			c.String(http.StatusBadRequest, "Missing 'search_keywords' query param")
			return
		}
	}

	periodHours := req.PeriodHours
	if periodHours <= 0 {
		periodHours = h.config.SearchPeriodHours
		if periodHours <= 0 {
			c.String(http.StatusBadRequest, "Missing 'period_hours' query param")
			return
		}
	}

	filterBy := req.FilterBy
	if len(filterBy) == 0 {
		filterBy = strings.Split(h.config.SearchFilter, ",")
	}

	scraper := scraper.NewScraper(h.db)
	_, err := scraper.ScrapeLinkedInJobs(c.Request.Context(), searchBy, filterBy, time.Duration(periodHours)*time.Hour)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("err: %v", err)
	c.Status(http.StatusOK)
}
