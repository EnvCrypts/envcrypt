-- +goose Up
ALTER TABLE projects ADD COLUMN prk_version INT NOT NULL DEFAULT 1;

-- +goose Down
ALTER TABLE projects DROP COLUMN prk_version;
