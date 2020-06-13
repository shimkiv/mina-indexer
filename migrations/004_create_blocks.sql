-- +goose Up
CREATE TABLE blocks (
  id                  SERIAL NOT NULL,
  app_version         TEXT,
  height              INTEGER NOT NULL,
  time                TIMESTAMP WITH TIME ZONE NOT NULL,
  hash                TEXT NOT NULL,
  parent_hash         TEXT NOT NULL,
  ledger_hash         TEXT NOT NULL,
  creator             TEXT NOT NULL,
  coinbase            DECIMAL(65, 0),
  total_currency      DECIMAL(65, 0),
  epoch               DOUBLE PRECISION,
  slot                INTEGER,
  transactions_count  INTEGER,
  fee_transfers_count INTEGER,
  snarkers_count      INTEGER,
  snark_jobs_count    INTEGER,
  snark_jobs_fees     DECIMAL(65, 0),
  created_at          TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at          TIMESTAMP WITH TIME ZONE NOT NULL,

  PRIMARY KEY (id)
);

CREATE INDEX idx_blocks_height
  ON blocks (height);

CREATE INDEX idx_blocks_hash
  ON blocks (hash);

CREATE INDEX idx_blocks_creator
  ON blocks (creator);

-- +goose Down
DROP TABLE IF EXISTS blocks;
