-- +goose Up
ALTER TABLE jobs 
ADD COLUMN date_posted TEXT;

-- +goose Down
ALTER TABLE jobs 
DROP COLUMN date_posted;