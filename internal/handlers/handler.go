package handlers

import (
	"github.com/tangerinefrog/GoScout/internal/data"
	"github.com/tangerinefrog/GoScout/internal/data/models"
	"github.com/tangerinefrog/GoScout/internal/state"
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
