-- +goose Up
ALTER TABLE jobs ADD COLUMN grade INT;
ALTER TABLE jobs ADD COLUMN grade_reasoning TEXT;

-- +goose Down
ALTER TABLE jobs DROP COLUMN grade;
ALTER TABLE jobs DROP COLUMN grade_reasoning;