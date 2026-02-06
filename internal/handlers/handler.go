package handlers

import (
	"github.com/tangerinefrog/GoScout/internal/data/models"
	"github.com/tangerinefrog/GoScout/internal/data/repositories"
	"github.com/tangerinefrog/GoScout/internal/state"
)

type handler struct {
	jobsRepository   *repositories.JobsRepository
	configRepository *repositories.ConfigRepository
	gradeState       *state.GradingState
	config           *models.Config
}

func NewHandler(jobsRepository *repositories.JobsRepository, configRepository *repositories.ConfigRepository) *handler {
	return &handler{
		jobsRepository:   jobsRepository,
		configRepository: configRepository,
		gradeState:       state.NewGradingLock(),
	}
}
