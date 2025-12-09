-- +goose Up
CREATE TABLE config (
    id INT PRIMARY KEY,
    search_query TEXT,
	search_filter TEXT,
	search_period_hours INT,
	grading_profile TEXT
);

-- +goose Down
DROP TABLE config;