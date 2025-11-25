package llm

import "job-scraper/internal/data"

type jobGrader struct {
	db *data.DB
}

type grade int

type GradeResult struct {
	Grade     grade
	Reasoning string
}

func NewJobGrader(db *data.DB) *jobGrader {
	return &jobGrader{
		db: db,
	}
}

func (jg *jobGrader) Grade(jobDescr string) (GradeResult, error) {
	return GradeResult{}, nil
}
