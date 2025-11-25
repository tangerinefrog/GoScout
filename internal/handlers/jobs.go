package handlers

import (
	"job-scraper/internal/data/repositories"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type JobResponse struct {
	Id            string    `json:"id"`
	Title         string    `json:"title"`
	Url           string    `json:"url"`
	Description   string    `json:"description"`
	DatePosted    time.Time `json:"date_posted"`
	Company       string    `json:"company"`
	Location      string    `json:"location"`
	NumApplicants string    `json:"numApplicants"`
	Status        string    `json:"status"`
}

func (h *handler) jobHandler(c *gin.Context) {
	jobId := c.Param("jobId")
	if jobId == "" {
		c.Status(http.StatusNotFound)
		return
	}

	jobsRepo := repositories.NewJobsRepo(h.db)
	job, err := jobsRepo.GetByID(jobId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := JobResponse{
		Id:            job.Id,
		Title:         job.Title,
		Url:           job.Url,
		Description:   job.Description,
		DatePosted:    job.DatePosted,
		Company:       job.Company,
		Location:      job.Location,
		NumApplicants: job.NumApplicants,
		Status:        string(job.Status),
	}

	c.JSON(http.StatusOK, resp)
}

func (h *handler) jobsHandler(c *gin.Context) {
	descrFlag := strings.ToLower(c.Query("includeDescr"))
	includeDescr := descrFlag == "true"

	jobsRepo := repositories.NewJobsRepo(h.db)
	jobs, err := jobsRepo.List()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := make([]JobResponse, len(jobs))

	for i, job := range jobs {
		result[i] = JobResponse{
			Id:            job.Id,
			Title:         job.Title,
			Url:           job.Url,
			DatePosted:    job.DatePosted,
			Company:       job.Company,
			Location:      job.Location,
			NumApplicants: job.NumApplicants,
			Status:        string(job.Status),
		}
		if includeDescr {
			result[i].Description = job.Description
		}
	}

	c.JSON(http.StatusOK, result)
}
