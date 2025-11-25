package handlers

import (
	"fmt"
	"job-scraper/internal/data/models"
	"job-scraper/internal/data/repositories"
	"job-scraper/internal/services/llm"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *handler) gradeHandler(c *gin.Context) {
	if h.gradeStatus.IsGrading() {
		c.Status(400)
		return
	}

	h.gradeStatus.Lock()

	jobsRepo := repositories.NewJobsRepo(h.db)
	ungradedJobs, err := jobsRepo.ListByStatus(models.JobStatusCreated)
	if err != nil {
		h.gradeStatus.Unlock()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	grader := llm.NewJobGrader(h.db)

	go func() {
		defer h.gradeStatus.Unlock()
		start := time.Now()
		for i, job := range ungradedJobs {
			status := fmt.Sprintf("grading job '%s'... jobs graded this cycle: %d", job.Id, i)
			h.gradeStatus.SetStatus(status)

			time.Sleep(10 * time.Second)
			_, err := grader.Grade(job.Description)
			if err != nil {
				status := fmt.Sprintf("error during grading job '%s': %v", job.Id, err)
				h.gradeStatus.SetStatus(status)
				return
			}
		}

		duration := time.Since(start)
		timestamp := time.Now().Format("15:04:05-07:00")

		status := fmt.Sprintf("grading completed at %s; run time: %s; total jobs graded: %d", timestamp, duration, len(ungradedJobs))
		h.gradeStatus.SetStatus(status)
	}()

	c.Status(200)
}

func (h *handler) gradeStatusHandler(c *gin.Context) {
	status := h.gradeStatus.GetStatus()

	c.JSON(http.StatusOK, gin.H{"status": status})
}
