package handlers

import (
	"github.com/gin-gonic/gin"
)

func (h *handler) SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	api.GET("/jobs", h.jobsHandler)
	api.POST("/scrape", h.scrapeHandler)
	api.GET("/jobs/:jobId", h.jobHandler)
	api.GET("/export", h.exportHandler)
	api.POST("/grade", h.gradeAllHandler)
	api.POST("/grade/:jobId", h.gradeJobHandler)
	api.GET("/grade/status", h.gradeStatusHandler)
}
