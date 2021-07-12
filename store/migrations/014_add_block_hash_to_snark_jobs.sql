-- +goose Up
ALTER TABLE snark_jobs ADD COLUMN block_hash TEXT;

CREATE INDEX idx_snark_jobs_block_hash
  ON snark_jobs(block_hash);

-- +goose Down
ALTER TABLE snark_jobs DROP COLUMN block_hash;

DROP INDEX IF EXISTS idx_snark_jobs_block_hash;