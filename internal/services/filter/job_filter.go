package filter

import (
	"context"
	"job-scraper/internal/data/models"
	"log"
	"strings"
)

type jobFilter struct {
	jRepo    JobsRepository
	keywords []string
}

type JobsRepository interface {
	GetByTitleAndCompany(ctx context.Context, title, company string) (*models.Job, error)
}

func NewJobFilter(jRepo JobsRepository, keywords []string) *jobFilter {
	return &jobFilter{
		jRepo:    jRepo,
		keywords: keywords,
	}
}

func (f *jobFilter) Filter(ctx context.Context, j models.Job) bool {
	if !containsKeywords(f.keywords, j.Description+j.Title) {
		return false
	}
	if isDuplicate(ctx, j.Title, j.Company, f.jRepo) {
		return false
	}

	return true
}

func containsKeywords(keywords []string, text string) bool {
	text = strings.ToLower(strings.TrimSpace(text))
	for _, keyword := range keywords {
		keyword = strings.ToLower(strings.TrimSpace(keyword))
		if keyword == "" {
			continue
		}
		matched := strings.Contains(text, keyword)
		if matched {
			return true
		}
	}

	return false
}

func isDuplicate(ctx context.Context, title, company string, jRepo JobsRepository) bool {
	job, err := jRepo.GetByTitleAndCompany(ctx, title, company)
	if err != nil {
		log.Printf("Error during getting a job by title: %v", err)
		return false
	}

	if job != nil {
		return true
	}

	return false
}
