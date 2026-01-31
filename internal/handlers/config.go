package handlers

import (
	"net/http"

	"github.com/tangerinefrog/GoScout/internal/data/models"
	"github.com/tangerinefrog/GoScout/internal/data/repositories"

	"github.com/gin-gonic/gin"
)

type Config struct {
	SearchQuery       string `json:"search_query"`
	SearchFilter      string `json:"search_filter"`
	SearchPeriodHours int    `json:"search_period_hours"`
	GradingProfile    string `json:"grading_profile"`
}

func (h *handler) configHandler(c *gin.Context) {
	configRepo := repositories.NewConfigRepo(h.db)
	config, err := configRepo.Get(c.Request.Context())
	if err != nil || config == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := Config{
		SearchQuery:       config.SearchQuery,
		SearchFilter:      config.SearchFilter,
		SearchPeriodHours: config.SearchPeriodHours,
		GradingProfile:    config.GradingProfile,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *handler) configUpdateHandler(c *gin.Context) {
	var req Config
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := models.Config{
		SearchQuery:       req.SearchQuery,
		SearchFilter:      req.SearchFilter,
		SearchPeriodHours: req.SearchPeriodHours,
		GradingProfile:    req.GradingProfile,
	}

	configRepo := repositories.NewConfigRepo(h.db)
	err := configRepo.Update(c.Request.Context(), &config)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
