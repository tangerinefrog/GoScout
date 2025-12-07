-- +goose Up
ALTER TABLE jobs ADD COLUMN is_invalid BOOLEAN DEFAULT 0;

-- +goose Down
ALTER TABLE jobs DROP COLUMN is_invalid;
