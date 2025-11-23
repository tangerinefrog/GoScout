-- +goose Up
CREATE TABLE jobs (
    id VARCHAR(20) PRIMARY KEY,
    title TEXT,
	url TEXT,
	description TEXT,
	timeAgo VARCHAR(20),
	company TEXT,
	location VARCHAR(100),
	numApplicants VARCHAR(10),
	status VARCHAR(20)
);

-- +goose Down
DROP TABLE jobs;