-- +goose Up
ALTER TABLE snark_jobs ADD COLUMN block_reference TEXT;

-- +goose Down
ALTER TABLE snark_jobs DROP COLUMN block_reference;