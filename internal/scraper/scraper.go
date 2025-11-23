package scraper

import (
	"errors"
	"fmt"
	"job-scraper/internal/fetcher"
	"job-scraper/internal/models"
	"job-scraper/internal/parser"
	"log"
	"net/url"
	"strconv"
	"strings"
)

const linkedInBaseUrl string = "https://www.linkedin.com/jobs-guest/jobs/api"

func ScrapeLinkedInJobs(jobTitle string) ([]models.JobPosition, error) {
	jobTitle = strings.TrimSpace(jobTitle)
	params := buildSearchQueryParams(jobTitle, 1)
	searchUrl := fmt.Sprintf("%s/seeMoreJobPostings/search?%s", linkedInBaseUrl, params)

	searchContent, err := fetcher.Fetch(searchUrl)
	if err != nil {
		return nil, err
	}

	jobIds, err := parser.ParseIdsFromSearch(searchContent)
	if err != nil {
		return nil, err
	}

	if len(jobIds) == 0 {
		return nil, errors.New("got no job IDs from search response")
	}

	res := make([]models.JobPosition, 0, len(jobIds))
	for _, jobId := range jobIds {
		jobPageUrl := fmt.Sprintf("%s/jobPosting/%s", linkedInBaseUrl, jobId)
		jobPostingContent, err := fetcher.Fetch(jobPageUrl)

		if err != nil {
			log.Printf("could not get job posting with id '%s' from '%s': %v\n", jobId, jobPageUrl, err)
			continue
		}

		job, err := parser.ParseJob(jobPostingContent, jobId)
		if err != nil {
			log.Printf("could not parse job posting with id '%s': %v\n", jobId, err)
			continue
		}
		res = append(res, job)
	}

	return res, nil
}

func buildSearchQueryParams(jobTitle string, page int) string {
	queryParams := url.Values{
		"keywords": {jobTitle},
		"location": {"Worldwide"},
		"start":    {strconv.Itoa((page - 1) * 10)},
	}

	return queryParams.Encode()
}
