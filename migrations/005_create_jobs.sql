-- +goose Up
CREATE TABLE IF NOT EXISTS jobs (
  id          BIGSERIAL NOT NULL,
  height      DOUBLE PRECISION NOT NULL,
  time        TIMESTAMP WITH TIME ZONE NOT NULL,
  prover      TEXT NOT NULL,
  fee         DECIMAL(65, 0) NOT NULL,
  works_count INTEGER NOT NULL DEFAULT 0,
  created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at  TIMESTAMP WITH TIME ZONE NOT NULL,

  PRIMARY KEY (id)
);

CREATE INDEX idx_jobs_time
  ON jobs(time);

CREATE INDEX idx_jobs_height
  ON jobs(height);

CREATE INDEX idx_jobs_prover
  ON jobs(prover);

-- +goose Down
DROP TABLE jobs;
