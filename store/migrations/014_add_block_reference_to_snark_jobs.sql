-- +goose Up
ALTER TABLE snark_jobs ADD COLUMN block_reference TEXT;

CREATE INDEX idx_snark_jobs_block_reference
  ON snark_jobs(block_reference);

-- +goose Down
ALTER TABLE snark_jobs DROP COLUMN block_reference;

DROP INDEX IF EXISTS idx_snark_jobs_block_reference;