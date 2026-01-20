package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"job-scraper/internal/data"
	"job-scraper/internal/data/models"
	"job-scraper/internal/data/sqltypes"
	"time"
)

type JobsRepo struct {
	db *data.DB
}

func NewJobsRepo(db *data.DB) *JobsRepo {
	return &JobsRepo{
		db: db,
	}
}

func (repo *JobsRepo) Add(ctx context.Context, job *models.Job) error {
	query := `
		INSERT INTO jobs
		(
			id,
			title,
			url,
			description,
			date_posted,
			company,
			location,
			num_applicants,
			status,
			note,
			is_invalid
		)
		VALUES
		(
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		);
	`

	_, err := repo.db.ExecContext(ctx, query,
		job.Id,
		job.Title,
		job.Url,
		job.Description,
		sqltypes.TimeToSqlFormat(job.DatePosted),
		job.Company,
		job.Location,
		job.NumApplicants,
		job.Status,
		job.Note,
		job.IsInvalid,
	)

	if err != nil {
		return fmt.Errorf("job insert failed: %w", err)
	}

	return nil
}

func (repo *JobsRepo) Update(ctx context.Context, job *models.Job) error {
	query := `
		UPDATE jobs
		SET 
			status = ?,
			grade = ?,
			grade_reasoning = ?,
			note = ?
		WHERE id = ?
	`

	_, err := repo.db.ExecContext(ctx, query,
		job.Status,
		job.Grade,
		job.GradeReasoning,
		job.Note,
		job.Id,
	)

	if err != nil {
		return fmt.Errorf("job update failed: %w", err)
	}

	return nil
}

func (repo *JobsRepo) GetByID(ctx context.Context, id string) (*models.Job, error) {
	query := `
		SELECT
			id,
			title,
			url,
			description,
			date_posted,
			company,
			location,
			num_applicants,
			status,
			grade,
			grade_reasoning,
			note,
			is_invalid
		FROM jobs
		WHERE id = ?;
	`

	rows, err := repo.db.QueryContext(ctx, query, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("jobs list query failed: %w", err)
	}

	defer rows.Close()

	hasJob := rows.Next()
	if !hasJob {
		return nil, nil
	}
	job := &models.Job{}
	var datePostedString string
	err = rows.Scan(
		&job.Id,
		&job.Title,
		&job.Url,
		&job.Description,
		&datePostedString,
		&job.Company,
		&job.Location,
		&job.NumApplicants,
		&job.Status,
		&job.Grade,
		&job.GradeReasoning,
		&job.Note,
		&job.IsInvalid,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan job row: %w", err)
	}

	job.DatePosted, _ = sqltypes.ParseTimeFromSql(datePostedString)

	return job, nil
}

func (repo *JobsRepo) List(ctx context.Context,
	status *models.JobStatus, company *string, grade *int, minDate *time.Time, search *string) ([]*models.Job, error) {
	query := `
		SELECT
			id,
			title,
			url,
			description,
			date_posted,
			company,
			location,
			num_applicants,
			status,
			grade,
			grade_reasoning,
			note
		FROM jobs
		WHERE 
			is_invalid = 0 AND is_archived = 0
			AND (? IS NULL OR ? = status)
			AND (? IS NULL OR ? = company) 
			AND (? IS NULL OR grade > ?) 
			AND (? IS NULL OR date_posted > ?)
			AND (? IS NULL OR description LIKE '%' || ? || '%')
		ORDER BY date_posted DESC;
	`
	var dateFilter *string
	if minDate != nil {
		d := sqltypes.TimeToSqlFormat(*minDate)
		dateFilter = &d
	}

	rows, err := repo.db.QueryContext(ctx, query,
		status, status,
		company, company,
		grade, grade,
		dateFilter, dateFilter,
		search, search,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("jobs list query failed: %w", err)
	}

	defer rows.Close()
	var jobs []*models.Job

	for rows.Next() {
		job := &models.Job{}
		var datePostedString string

		err := rows.Scan(
			&job.Id,
			&job.Title,
			&job.Url,
			&job.Description,
			&datePostedString,
			&job.Company,
			&job.Location,
			&job.NumApplicants,
			&job.Status,
			&job.Grade,
			&job.GradeReasoning,
			&job.Note,
		)
		if err != nil {
			return jobs, fmt.Errorf("failed to scan job row: %w", err)
		}
		job.DatePosted, _ = sqltypes.ParseTimeFromSql(datePostedString)

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (repo *JobsRepo) Archive(ctx context.Context, id string) error {
	query := `
		UPDATE jobs
		SET 
			is_archived = 1,
			description = ''
		WHERE id = ?
	`

	_, err := repo.db.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("job update failed: %w", err)
	}

	return nil
}
