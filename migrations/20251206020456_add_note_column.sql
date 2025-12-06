-- +goose Up
ALTER TABLE jobs ADD COLUMN note TEXT NOT NULL DEFAULT('');

-- +goose Down
ALTER TABLE jobs DROP COLUMN note;
