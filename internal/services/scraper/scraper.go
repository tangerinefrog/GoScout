package scraper

import (
	"context"
	"errors"
	"fmt"
	"job-scraper/internal/data"
	"job-scraper/internal/data/models"
	"job-scraper/internal/data/repositories"
	"job-scraper/internal/services/fetcher"
	"job-scraper/internal/services/parser"
	"job-scraper/internal/services/validator"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const linkedInBaseUrl string = "https://www.linkedin.com/jobs-guest/jobs/api"

type Scraper struct {
	db *data.DB
}

func NewScraper(db *data.DB) *Scraper {
	return &Scraper{
		db: db,
	}
}

func (s *Scraper) ScrapeLinkedInJobs(ctx context.Context, keyword string, filterKeywords []string, timeWindow time.Duration) ([]models.Job, error) {
	keyword = strings.TrimSpace(keyword)

	jobIds, err := getJobsFromSearch(ctx, keyword, timeWindow)
	if err != nil {
		return nil, err
	}

	validator := validator.NewKeywordValidator()
	jRepo := repositories.NewJobsRepo(s.db)
	res := make([]models.Job, 0, len(jobIds))
	for _, jobId := range jobIds {
		dbJob, err := jRepo.GetByID(ctx, jobId)
		if err != nil {
			log.Printf("could not get job with id '%s' from database: %v\n", jobId, err)
			continue
		}
		if dbJob != nil {
			res = append(res, *dbJob)
			continue
		}

		jobPageUrl := fmt.Sprintf("%s/jobPosting/%s", linkedInBaseUrl, jobId)
		jobPostingContent, err := fetcher.FetchWithRetry(ctx, jobPageUrl, 5)

		if err != nil {
			log.Printf("could not get job with id '%s' from '%s': %v\n", jobId, jobPageUrl, err)
			continue
		}

		job, err := parser.ParseJob(jobPostingContent, jobId)
		if err != nil {
			log.Printf("could not parse job with id '%s': %v\n", jobId, err)
			continue
		}
		job.Status = models.JobStatusCreated

		valid := validator.ValidateKeywords(filterKeywords, job.Description+job.Title)
		if !valid {
			job.IsInvalid = true
		}

		err = jRepo.Add(ctx, &job)
		if err != nil {
			log.Printf("could not save job with id '%s' from '%s' to database: %v\n", jobId, jobPageUrl, err)
			continue
		}

		res = append(res, job)
	}

	return res, nil
}

func getJobsFromSearch(ctx context.Context, keywords string, timeWindow time.Duration) ([]string, error) {
	var ids []string

	page := 1
	for {
		params := buildSearchQueryParams(keywords, page, timeWindow)
		searchUrl := fmt.Sprintf("%s/seeMoreJobPostings/search?%s", linkedInBaseUrl, params)

		searchContent, err := fetcher.Fetch(ctx, searchUrl)
		if err != nil {
			if !errors.Is(err, fetcher.ErrorUnsuccessfulStatusCode) {
				return nil, err
			}
		}
		jobIds, err := parser.ParseIdsFromSearch(searchContent)
		if err != nil {
			return nil, err
		}
		if len(jobIds) == 0 {
			break
		}

		ids = append(ids, jobIds...)

		page++
	}

	return ids, nil
}

func buildSearchQueryParams(keywords string, page int, timeWindow time.Duration) string {
	if timeWindow == 0 {
		timeWindow = 24 * time.Hour
	}

	queryParams := url.Values{
		"keywords": {keywords},
		//todo: hardcoded, can change to a country later
		"location": {"Worldwide"},
		//get jobs posts for the last N seconds
		"f_TPR": {fmt.Sprintf("r%.0f", timeWindow.Seconds())},
		//remote work
		"f_WT":  {"2"},
		"start": {strconv.Itoa((page - 1) * 25)},
	}

	return queryParams.Encode()
}
