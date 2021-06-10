-- +goose Up
ALTER TABLE transactions ADD COLUMN canonical BOOLEAN;

-- +goose Down
ALTER TABLE transactions DROP COLUMN canonical;