-- +goose Up
CREATE TABLE snark_jobs (
  id          SERIAL NOT NULL,
  height      CHAIN_HEIGHT,
  time        CHAIN_TIME,
  prover      TEXT NOT NULL,
  fee         CHAIN_CURRENCY,
  works_count INTEGER NOT NULL,
  created_at  CHAIN_TIME,

  PRIMARY KEY (id)
);

CREATE INDEX idx_snark_jobs_time
  ON snark_jobs(time);

CREATE INDEX idx_snark_jobs_height
  ON snark_jobs(height);

CREATE INDEX idx_snark_jobs_prover
  ON snark_jobs(prover);

-- +goose Down
DROP TABLE snark_jobs;
