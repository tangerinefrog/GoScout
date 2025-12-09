package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"job-scraper/internal/data"
	"job-scraper/internal/data/models"
)

type ConfigRepo struct {
	db *data.DB
}

func NewConfigRepo(db *data.DB) *ConfigRepo {
	return &ConfigRepo{
		db: db,
	}
}

func (repo *ConfigRepo) Init(ctx context.Context) error {
	selectQuery := `
		SELECT 1 FROM config WHERE id = 1;
	`

	rows, err := repo.db.QueryContext(ctx, selectQuery)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("config select query failed: %w", err)
	}

	defer rows.Close()
	if rows.Next() {
		return nil
	}

	insertQuery := `
		INSERT INTO config
		(
			id, 
			search_query, 
			search_filter,
			search_period_hours,
			grading_profile
		)
		VALUES
		(
			1,
			'',
			'',
			1,
			''
		);
	`

	_, err = repo.db.ExecContext(ctx, insertQuery)

	if err != nil {
		return fmt.Errorf("config insert query failed: %w", err)
	}

	return nil
}

func (repo *ConfigRepo) Update(ctx context.Context, config *models.Config) error {
	query := `
		UPDATE config
		SET 
			search_query = ?,
			search_filter = ?,
			search_period_hours = ?,
			grading_profile = ?
		WHERE id = 1
	`

	_, err := repo.db.ExecContext(ctx, query,
		config.SearchQuery,
		config.SearchFilter,
		config.SearchPeriodHours,
		config.GradingProfile,
	)

	if err != nil {
		return fmt.Errorf("config update query failed: %w", err)
	}

	return nil
}

func (repo *ConfigRepo) Get(ctx context.Context) (*models.Config, error) {
	query := `
		SELECT
			search_query, 
			search_filter,
			search_period_hours,
			grading_profile
		FROM config
		WHERE id = 1;
	`

	rows, err := repo.db.QueryContext(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("config select query failed: %w", err)
	}

	defer rows.Close()

	hasConfig := rows.Next()
	if !hasConfig {
		return nil, nil
	}
	config := &models.Config{}

	err = rows.Scan(
		&config.SearchQuery,
		&config.SearchFilter,
		&config.SearchPeriodHours,
		&config.GradingProfile,
	)
	if err != nil {
		return nil, fmt.Errorf("config scan from db query failed: %w", err)
	}

	return config, nil
}
