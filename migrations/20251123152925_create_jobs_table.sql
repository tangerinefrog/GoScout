-- +goose Up
CREATE TABLE jobs (
    id VARCHAR(20) PRIMARY KEY,
    title TEXT,
	url TEXT,
	description TEXT,
	company TEXT,
	location VARCHAR(100),
	num_applicants VARCHAR(10),
	status VARCHAR(20)
);

-- +goose Down
DROP TABLE jobs;