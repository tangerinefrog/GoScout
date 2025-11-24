package handlers

import (
	"job-scraper/internal/data"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	db *data.DB
}

func NewHandler(db *data.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) SetupRoutes(r *gin.Engine) {
	r.POST("api/scrape", h.scrapeHandler)
	r.GET("api/jobs", h.jobsHandler)
	r.GET("api/jobs/:jobId", h.jobHandler)
}
