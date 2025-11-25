package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"job-scraper/internal/data"
	"job-scraper/internal/data/models"
	"job-scraper/internal/data/sqltypes"
)

type JobsRepo struct {
	db *data.DB
}

func NewJobsRepo(db *data.DB) *JobsRepo {
	return &JobsRepo{
		db: db,
	}
}

func (repo *JobsRepo) Add(job *models.Job) error {
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
			status
		)
		VALUES
		(
			?, ?, ?, ?, ?, ?, ?, ?, ?
		);
	`

	_, err := repo.db.ExecContext(context.TODO(), query,
		job.Id,
		job.Title,
		job.Url,
		job.Description,
		sqltypes.TimeToSqlFormat(job.DatePosted),
		job.Company,
		job.Location,
		job.NumApplicants,
		job.Status,
	)

	if err != nil {
		return fmt.Errorf("job insert failed: %w", err)
	}

	return nil
}

func (repo *JobsRepo) Update(job *models.Job) error {
	query := `
		UPDATE jobs
		SET 
			status = ?,
			grade = ?,
			grade_reasoning = ?
		WHERE id = ?
	`

	_, err := repo.db.ExecContext(context.TODO(), query,
		job.Status,
		job.Grade,
		job.GradeReasoning,
		job.Id,
	)

	if err != nil {
		return fmt.Errorf("job update failed: %w", err)
	}

	return nil
}

func (repo *JobsRepo) GetByID(id string) (*models.Job, error) {
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
			grade_reasoning
		FROM jobs
		WHERE id = ?;
	`

	rows, err := repo.db.QueryContext(context.TODO(), query, id)

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
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan job row: %w", err)
	}

	job.DatePosted = sqltypes.ParseTimeFromSql(datePostedString)

	return job, nil
}

func (repo *JobsRepo) List() ([]*models.Job, error) {
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
			grade_reasoning
		FROM jobs;
	`

	rows, err := repo.db.QueryContext(context.TODO(), query)
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
		)
		if err != nil {
			return jobs, fmt.Errorf("failed to scan job row: %w", err)
		}
		job.DatePosted = sqltypes.ParseTimeFromSql(datePostedString)

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (repo *JobsRepo) ListByStatus(status models.JobStatus) ([]*models.Job, error) {
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
			grade_reasoning
		FROM jobs
		WHERE status = ?;
	`

	rows, err := repo.db.QueryContext(context.TODO(), query, status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("jobs list by status query failed: %w", err)
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
		)
		if err != nil {
			return jobs, fmt.Errorf("failed to scan job row: %w", err)
		}
		job.DatePosted = sqltypes.ParseTimeFromSql(datePostedString)

		jobs = append(jobs, job)
	}

	return jobs, nil
}
