package handlers

import (
	"job-scraper/internal/data"
	"job-scraper/internal/state"
)

type handler struct {
	db          *data.DB
	gradeStatus *state.GradingState
}

func NewHandler(db *data.DB) *handler {
	return &handler{
		db:          db,
		gradeStatus: state.NewGradingLock(),
	}
}
