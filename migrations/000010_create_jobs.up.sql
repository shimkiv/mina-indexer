CREATE TABLE IF NOT EXISTS jobs (
  id         BIGSERIAL                NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
  height     DOUBLE PRECISION         NOT NULL,
  time       TIMESTAMP WITH TIME ZONE NOT NULL,
  prover     TEXT                     NOT NULL,
  fee        DECIMAL(65, 0)           NOT NULL,

  PRIMARY KEY (id, time)
);

-- Hypertable
SELECT create_hypertable('jobs', 'time', if_not_exists => TRUE);

-- Indexes
CREATE INDEX idx_jobs_prover      ON jobs(prover);
CREATE INDEX idx_jobs_height      ON jobs(height);
CREATE INDEX idx_jobs_height_time ON jobs(height, time);