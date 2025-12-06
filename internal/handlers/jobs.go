package handlers

import (
	"job-scraper/internal/data/models"
	"job-scraper/internal/data/repositories"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type JobResponse struct {
	Id             string    `json:"id"`
	Title          string    `json:"title"`
	Url            string    `json:"url"`
	Description    string    `json:"description"`
	DatePosted     time.Time `json:"date_posted"`
	Company        string    `json:"company"`
	Location       string    `json:"location"`
	NumApplicants  string    `json:"num_applicants"`
	Status         string    `json:"status"`
	Grade          *int      `json:"grade"`
	GradeReasoning *string   `json:"grade_reasoning"`
	Note           string    `json:"note"`
}

type JobUpdateRequest struct {
	Status string  `json:"status"`
	Grade  int     `json:"grade"`
	Note   *string `json:"note"`
}

func (h *handler) jobHandler(c *gin.Context) {
	jobId := c.Param("jobId")
	if jobId == "" {
		c.Status(http.StatusNotFound)
		return
	}

	jobsRepo := repositories.NewJobsRepo(h.db)
	job, err := jobsRepo.GetByID(c.Request.Context(), jobId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if job == nil {
		c.Status(http.StatusNotFound)
		return
	}

	resp := JobResponse{
		Id:             job.Id,
		Title:          job.Title,
		Url:            job.Url,
		Description:    job.Description,
		DatePosted:     job.DatePosted,
		Company:        job.Company,
		Location:       job.Location,
		NumApplicants:  job.NumApplicants,
		Status:         string(job.Status),
		Grade:          job.Grade,
		GradeReasoning: job.GradeReasoning,
		Note:           job.Note,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *handler) jobsHandler(c *gin.Context) {
	descrFlag := strings.ToLower(c.Query("include_descr"))
	statusParam := strings.ToLower(strings.TrimSpace(c.Query("status")))
	companyParam := strings.TrimSpace(c.Query("company"))
	gradeParam := strings.TrimSpace(c.Query("grade_gt"))
	dateParam := strings.TrimSpace(c.Query("date_gt"))

	includeDescr := descrFlag == "true"
	var statusFilter *models.JobStatus
	if statusParam != "" {
		s := models.JobStatus(statusParam)
		statusFilter = &s
	}

	var companyFilter *string
	if companyParam != "" {
		companyFilter = &companyParam
	}

	var gradeFilter *int
	grade, err := strconv.Atoi(gradeParam)
	if err == nil && grade > 0 && grade < 5 {
		gradeFilter = &grade
	}

	var dateFilter *time.Time
	date, err := time.Parse("2006-01-02", dateParam)
	if err == nil {
		dateFilter = &date
	}

	jobsRepo := repositories.NewJobsRepo(h.db)
	jobs, err := jobsRepo.List(c.Request.Context(), statusFilter, companyFilter, gradeFilter, dateFilter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := make([]JobResponse, len(jobs))

	for i, job := range jobs {
		result[i] = JobResponse{
			Id:             job.Id,
			Title:          job.Title,
			Url:            job.Url,
			DatePosted:     job.DatePosted,
			Company:        job.Company,
			Location:       job.Location,
			NumApplicants:  job.NumApplicants,
			Status:         string(job.Status),
			Grade:          job.Grade,
			GradeReasoning: job.GradeReasoning,
			Note:           job.Note,
		}
		if includeDescr {
			result[i].Description = job.Description
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *handler) updateJobHandler(c *gin.Context) {
	jobId := c.Param("jobId")
	if jobId == "" {
		c.Status(http.StatusNotFound)
		return
	}

	jobsRepo := repositories.NewJobsRepo(h.db)

	job, err := jobsRepo.GetByID(c.Request.Context(), jobId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if job == nil {
		c.Status(http.StatusNotFound)
		return
	}

	var req JobUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Status != "" {
		job.Status = models.JobStatus(strings.ToLower(req.Status))
	}
	if req.Grade > 0 && req.Grade < 6 {
		job.Grade = &req.Grade
	}
	if req.Note != nil {
		job.Note = *req.Note
	}

	err = jobsRepo.Update(c.Request.Context(), job)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := JobResponse{
		Id:             job.Id,
		Title:          job.Title,
		Url:            job.Url,
		Description:    job.Description,
		DatePosted:     job.DatePosted,
		Company:        job.Company,
		Location:       job.Location,
		NumApplicants:  job.NumApplicants,
		Status:         string(job.Status),
		Grade:          job.Grade,
		GradeReasoning: job.GradeReasoning,
		Note:           job.Note,
	}

	c.JSON(http.StatusOK, resp)
}
