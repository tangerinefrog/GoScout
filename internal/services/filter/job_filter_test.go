package filter

import (
	"context"
	"errors"
	"job-scraper/internal/data/models"
	"testing"
)

type mockJobRepo struct {
	getByTitleAndCompanyFunc func(ctx context.Context, title, company string) (*models.Job, error)
}

func (m *mockJobRepo) GetByTitleAndCompany(ctx context.Context, title, company string) (*models.Job, error) {
	if m.getByTitleAndCompanyFunc != nil {
		return m.getByTitleAndCompanyFunc(ctx, title, company)
	}
	return nil, nil
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name     string
		keywords []string
		job      models.Job
		mockRepo *mockJobRepo
		want     bool
	}{
		{
			name:     "valid job",
			keywords: []string{"golang", "backend"},
			job: models.Job{
				Title:       "Senior Golang Developer",
				Company:     "AnyCorp",
				Description: "AnyDescription",
			},
			mockRepo: &mockJobRepo{
				getByTitleAndCompanyFunc: func(ctx context.Context, title, company string) (*models.Job, error) {
					return nil, nil
				},
			},
			want: true,
		},
		{
			name:     "keyword match fail",
			keywords: []string{"python"},
			job: models.Job{
				Title:       "Senior Golang Developer",
				Company:     "AnyCorp",
				Description: "AnyDescription",
			},
			mockRepo: &mockJobRepo{
				getByTitleAndCompanyFunc: func(ctx context.Context, title, company string) (*models.Job, error) {
					return nil, nil
				},
			},
			want: false,
		},
		{
			name:     "duplicated job",
			keywords: []string{"golang"},
			job: models.Job{
				Title:       "Golang Developer",
				Company:     "AnyCorp",
				Description: "AnyDescription",
			},
			mockRepo: &mockJobRepo{
				getByTitleAndCompanyFunc: func(ctx context.Context, title, company string) (*models.Job, error) {
					return &models.Job{
						Id: "123",
					}, nil
				},
			},
			want: false,
		},
		{
			name:     "repository returns error",
			keywords: []string{"golang"},
			job: models.Job{
				Title:       "Golang Developer",
				Company:     "AnyCorp",
				Description: "AnyDescription",
			},
			mockRepo: &mockJobRepo{
				getByTitleAndCompanyFunc: func(ctx context.Context, title, company string) (*models.Job, error) {
					return nil, errors.New("database error")
				},
			},
			want: true,
		},
		{
			name:     "keyword matches in description only",
			keywords: []string{"golang"},
			job: models.Job{
				Title:       "DevOps Engineer",
				Company:     "AnyCorp",
				Description: "Experience with Golang",
			},
			mockRepo: &mockJobRepo{
				getByTitleAndCompanyFunc: func(ctx context.Context, title, company string) (*models.Job, error) {
					return nil, nil
				},
			},
			want: true,
		},
		{
			name:     "keyword matches in title only",
			keywords: []string{"senior"},
			job: models.Job{
				Title:       "Senior Developer",
				Company:     "AnyCorp",
				Description: "AnyDescription",
			},
			mockRepo: &mockJobRepo{
				getByTitleAndCompanyFunc: func(ctx context.Context, title, company string) (*models.Job, error) {
					return nil, nil
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewJobFilter(tt.mockRepo, tt.keywords)
			got := filter.Filter(context.Background(), tt.job)

			if got != tt.want {
				t.Errorf("Filter() got '%v', expected '%v'", got, tt.want)
			}
		})
	}
}

func TestContainsKeywords(t *testing.T) {
	tests := []struct {
		name     string
		keywords []string
		content  string
		want     bool
	}{
		{
			name:     "matches single keyword",
			keywords: []string{"golang"},
			content:  "Senior golang developer",
			want:     true,
		},
		{
			name:     "no match",
			keywords: []string{"golang"},
			content:  "Looking for a middle frontend developer with 3+ years of experience",
			want:     false,
		},
		{
			name:     "empty keywords",
			keywords: []string{},
			content:  "Sample job posting text",
			want:     false,
		},
		{
			name:     "matches at least one keyword",
			keywords: []string{"golang", "middle"},
			content:  "Looking for a middle backend developer with 3+ years of experience",
			want:     true,
		},
		{
			name:     "case insensitive match",
			keywords: []string{"gOlAnG"},
			content:  "Golang developer",
			want:     true,
		},
		{
			name:     "partial word match",
			keywords: []string{"Go"},
			content:  "Golang developer",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsKeywords(tt.keywords, tt.content)
			if got != tt.want {
				t.Errorf("containsKeywords() got '%v', expected '%v'", got, tt.want)
			}
		})
	}
}

func TestIsDuplicate(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		company  string
		mockRepo *mockJobRepo
		want     bool
	}{
		{
			name:    "no duplicates",
			title:   "Golang Developer",
			company: "AnyCorp",
			mockRepo: &mockJobRepo{
				getByTitleAndCompanyFunc: func(ctx context.Context, title, company string) (*models.Job, error) {
					return nil, nil
				},
			},
			want: false,
		},
		{
			name:    "has duplicate",
			title:   "Golang Developer",
			company: "AnyCorp",
			mockRepo: &mockJobRepo{
				getByTitleAndCompanyFunc: func(ctx context.Context, title, company string) (*models.Job, error) {
					return &models.Job{
						Id: "123",
					}, nil
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isDuplicate(context.Background(), tt.title, tt.company, tt.mockRepo)

			if got != tt.want {
				t.Errorf("isDuplicate() got '%v', expected '%v'", got, tt.want)
			}
		})
	}
}
