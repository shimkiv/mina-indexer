-- +goose Up
ALTER TABLE transactions ADD COLUMN canonical BOOLEAN NOT NULL;

-- +goose Down
ALTER TABLE transactions DROP COLUMN canonical;