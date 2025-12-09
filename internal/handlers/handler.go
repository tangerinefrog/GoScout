package handlers

import (
	"job-scraper/internal/data"
	"job-scraper/internal/data/models"
	"job-scraper/internal/state"
)

type handler struct {
	db         *data.DB
	gradeState *state.GradingState
	config     *models.Config
}

func NewHandler(db *data.DB) *handler {
	return &handler{
		db:         db,
		gradeState: state.NewGradingLock(),
	}
}
