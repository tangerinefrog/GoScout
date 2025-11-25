package handlers

import (
	"job-scraper/internal/data"

	"github.com/gin-gonic/gin"
)

type handler struct {
	db *data.DB
}

func NewHandler(db *data.DB) *handler {
	return &handler{
		db: db,
	}
}

func (h *handler) SetupRoutes(r *gin.Engine) {
	r.POST("api/scrape", h.scrapeHandler)
	r.GET("api/jobs", h.jobsHandler)
	r.GET("api/jobs/:jobId", h.jobHandler)
	r.GET("api/export", h.exportHandler)
}
