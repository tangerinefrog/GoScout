package handlers

import (
	"github.com/gin-gonic/gin"
)

func (h *handler) SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	api.POST("/scrape", h.scrapeHandler)

	api.GET("/jobs", h.jobsHandler)
	api.GET("/jobs/:jobId", h.jobHandler)
	api.PATCH("/jobs/:jobId", h.jobUpdateHandler)

	api.POST("/grade", h.gradeAllHandler)
	api.POST("/grade/:jobId", h.gradeJobHandler)
	api.GET("/grade/status", h.gradeStatusHandler)
	api.POST("/grade/stop", h.gradingStopHandler)

	api.GET("/export", h.exportHandler)

	api.GET("/config", h.configHandler)
	api.PUT("/config", h.configUpdateHandler)

	r.LoadHTMLFiles("web/index.html")
	r.Static("/static", "./web/static")
	r.GET("/", pageHandler)
}
