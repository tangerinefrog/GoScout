package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tangerinefrog/GoScout/internal/data/models"
	"github.com/tangerinefrog/GoScout/internal/data/repositories"
	"github.com/tangerinefrog/GoScout/internal/services/llm"

	"github.com/gin-gonic/gin"
)

type GradeRequest struct {
	Requirements string `json:"requirements"`
}

func (h *handler) gradeAllHandler(c *gin.Context) {
	var req GradeRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requirements := req.Requirements
	if requirements == "" {
		requirements = h.config.GradingProfile
		if requirements == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing 'requirements' body param"})
			return
		}
	}

	if h.gradeState.IsGrading() {
		status := h.gradeState.GetStatus()
		errMsg := fmt.Sprintf("another grading process is working, please wait. status: %s", status)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	h.gradeState.Lock()

	jobsRepo := repositories.NewJobsRepo(h.db)

	status := models.JobStatusCreated
	ungradedJobs, err := jobsRepo.List(c.Request.Context(), &status, nil, nil, nil, nil)
	if err != nil {
		h.gradeState.Unlock()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//grading takes a long time, let's restrict grading to N jobs per request
	const batchLen = 50
	var jobsBatch []*models.Job
	if len(ungradedJobs) > batchLen {
		jobsBatch = ungradedJobs[:batchLen]
	} else {
		jobsBatch = ungradedJobs
	}

	grader := llm.NewJobGrader(h.db)

	go func() {
		defer h.gradeState.Unlock()
		start := time.Now()
		for i, job := range jobsBatch {
			status := fmt.Sprintf("grading job '%s' (%d of %d)", job.Id, i+1, len(jobsBatch))
			h.gradeState.SetStatus(status)
			err := gradeJob(h.gradeState.Ctx, job, requirements, grader, jobsRepo)
			if err != nil {
				h.gradeState.SetStatus(err.Error())
				return
			}
		}

		duration := time.Since(start)
		timestamp := time.Now().Format("15:04:05-07:00")

		status := fmt.Sprintf("grading completed at %s; run time: %s; total jobs graded: %d", timestamp, duration, len(jobsBatch))
		h.gradeState.SetStatus(status)
	}()

	c.Status(200)
}

func (h *handler) gradeJobHandler(c *gin.Context) {
	jobId := c.Param("jobId")
	if jobId == "" {
		c.Status(http.StatusNotFound)
		return
	}

	var req GradeRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requirements := req.Requirements
	if requirements == "" {
		requirements = h.config.GradingProfile
		if requirements == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing 'requirements' body param"})
			return
		}
	}

	if h.gradeState.IsGrading() {
		status := h.gradeState.GetStatus()
		errMsg := fmt.Sprintf("another grading process is working, please wait. status: %s", status)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	h.gradeState.Lock()

	jobsRepo := repositories.NewJobsRepo(h.db)
	job, err := jobsRepo.GetByID(c.Request.Context(), jobId)
	if err != nil {
		h.gradeState.Unlock()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if job == nil {
		h.gradeState.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("no job with id '%s'", jobId)})
		return
	}

	grader := llm.NewJobGrader(h.db)

	go func() {
		defer h.gradeState.Unlock()
		start := time.Now()
		status := fmt.Sprintf("grading job '%s' (1 of 1)", job.Id)
		h.gradeState.SetStatus(status)

		err := gradeJob(h.gradeState.Ctx, job, requirements, grader, jobsRepo)
		if err != nil {
			h.gradeState.SetStatus(err.Error())
			return
		}

		duration := time.Since(start)
		timestamp := time.Now().Format("15:04:05-07:00")

		status = fmt.Sprintf("grading completed at %s; run time: %s", timestamp, duration)
		h.gradeState.SetStatus(status)
	}()

	c.Status(200)
}

func (h *handler) gradingStopHandler(c *gin.Context) {
	h.gradeState.Cancel()
}

func gradeJob(ctx context.Context, job *models.Job, requirements string, grader *llm.JobGrader,
	jobsRepo *repositories.JobsRepository) error {
	res, err := grader.Grade(ctx, requirements, job.Description)
	if err != nil {
		return fmt.Errorf("error during grading job '%s': %w", job.Id, err)
	}

	job.Status = models.JobStatusPending
	job.Grade = &res.Grade
	job.GradeReasoning = &res.Reasoning
	err = jobsRepo.Update(ctx, job)
	if err != nil {
		return fmt.Errorf("error during updating job '%s' grade: %w", job.Id, err)
	}

	return nil
}

func (h *handler) gradeStatusHandler(c *gin.Context) {
	status := h.gradeState.GetStatus()

	c.JSON(http.StatusOK, gin.H{"status": status})
}
