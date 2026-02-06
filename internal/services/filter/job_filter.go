package filter

import (
	"context"
	"log"
	"strings"

	"github.com/tangerinefrog/GoScout/internal/data/models"
)

type JobFilter struct {
	jRepo    jobsRepository
	keywords []string
}

type jobsRepository interface {
	GetByTitleAndCompany(ctx context.Context, title, company string) (*models.Job, error)
}

func NewJobFilter(jRepo jobsRepository, keywords []string) *JobFilter {
	return &JobFilter{
		jRepo:    jRepo,
		keywords: keywords,
	}
}

func (f *JobFilter) Filter(ctx context.Context, j models.Job) bool {
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

func isDuplicate(ctx context.Context, title, company string, jRepo jobsRepository) bool {
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
