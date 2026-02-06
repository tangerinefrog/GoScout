package scraper

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tangerinefrog/GoScout/internal/data"
	"github.com/tangerinefrog/GoScout/internal/data/models"
	"github.com/tangerinefrog/GoScout/internal/data/repositories"
	"github.com/tangerinefrog/GoScout/internal/services/fetcher"
	"github.com/tangerinefrog/GoScout/internal/services/filter"
	"github.com/tangerinefrog/GoScout/internal/services/parser"
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

func (s *Scraper) ScrapeLinkedInJobs(ctx context.Context, keyword string, filterKeywords []string, timeWindow time.Duration) error {
	keyword = strings.TrimSpace(keyword)

	jobIds, err := searchJobs(ctx, keyword, timeWindow)
	if err != nil {
		return err
	}


	jobsCh := make(chan string)
	go func() {
		for _, jobId := range jobIds {
			select {
			case jobsCh <- jobId:
			case <-ctx.Done():
				return
			}
		}

		close(jobsCh)
	}()

	jRepo := repositories.NewJobsRepo(s.db)
	jobFilter := filter.NewJobFilter(jRepo, filterKeywords)
	numWorkers := 3

	workerPool(ctx, jobsCh, numWorkers, jRepo, *jobFilter)

	return nil
}

func workerPool(ctx context.Context, jobsCh <-chan string, numWorkers int, r *repositories.JobsRepository, f filter.JobFilter) {
	wg := sync.WaitGroup{}
	for range numWorkers {
		wg.Go(func() {
			for {
				select {
				case id, ok := <-jobsCh:
					if !ok {
						return
					}
					scrapeAndSaveJob(ctx, id, r, f)
				case <-ctx.Done():
					return
				}
			}
		})
	}

	wg.Wait()
}

func scrapeAndSaveJob(ctx context.Context, id string, r *repositories.JobsRepository, f filter.JobFilter) {
	dbJob, err := r.GetByID(ctx, id)
	if err != nil {
		log.Printf("Error while getting a job with id '%s' from database: %v\n", id, err)
		return
	}
	if dbJob != nil {
		return
	}

	job, err := scrapeJob(ctx, id)
	if err != nil {
		log.Printf("Error while fetching a job: %v", err)
		return
	}

	job.Status = models.JobStatusCreated

	valid := f.Filter(ctx, job)
	if !valid {
		job.IsInvalid = true
		job.Description = ""
	}

	err = r.Add(ctx, &job)
	if err != nil {
		log.Printf("Could not save job with id '%s' to database: %v\n", id, err)
		return
	}
}

func scrapeJob(ctx context.Context, id string) (models.Job, error) {
	jobPageUrl := fmt.Sprintf("%s/jobPosting/%s", linkedInBaseUrl, id)
	jobPostingContent, err := fetcher.FetchWithRetry(ctx, jobPageUrl, 5)

	if err != nil {
		return models.Job{}, fmt.Errorf("could not get job with id '%s' from '%s': %w", id, jobPageUrl, err)
	}

	job, err := parser.ParseJob(jobPostingContent, id)
	if err != nil {
		return models.Job{}, fmt.Errorf("could not parse job with id '%s': %w", id, err)
	}

	return job, nil
}

func searchJobs(ctx context.Context, keywords string, timeWindow time.Duration) ([]string, error) {
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
