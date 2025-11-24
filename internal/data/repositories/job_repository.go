package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"job-scraper/internal/data"
	"job-scraper/internal/data/models"
	"job-scraper/internal/data/sqltypes"
)

type jobsRepo struct {
	db *data.DB
}

func NewJobsRepo(db *data.DB) *jobsRepo {
	return &jobsRepo{
		db: db,
	}
}

func (repo *jobsRepo) Add(job *models.Job) error {
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
		)
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

func (repo *jobsRepo) GetByID(id string) (*models.Job, error) {
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
			status
		FROM jobs
		WHERE id = ?
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
	var job models.Job
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
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan job row: %w", err)
	}

	job.DatePosted = sqltypes.ParseTimeFromSql(datePostedString)

	return &job, nil
}

func (repo *jobsRepo) List() ([]models.Job, error) {
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
			status
		FROM jobs
	`

	rows, err := repo.db.QueryContext(context.TODO(), query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("jobs list query failed: %w", err)
	}

	defer rows.Close()
	var jobs []models.Job

	for rows.Next() {
		var job models.Job
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
		)
		if err != nil {
			return jobs, fmt.Errorf("failed to scan job row: %w", err)
		}
		job.DatePosted = sqltypes.ParseTimeFromSql(datePostedString)

		jobs = append(jobs, job)
	}

	return jobs, nil
}
