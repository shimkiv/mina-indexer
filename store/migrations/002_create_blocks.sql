-- +goose Up
CREATE TABLE blocks (
  id                   SERIAL NOT NULL,
  canonical            BOOLEAN NOT NULL,
  height               CHAIN_HEIGHT,
  time                 CHAIN_TIME,
  hash                 TEXT NOT NULL,
  parent_hash          TEXT NOT NULL,
  ledger_hash          TEXT NOT NULL,
  snarked_ledger_hash  TEXT NOT NULL,
  creator              TEXT NOT NULL,
  coinbase             CHAIN_CURRENCY DEFAULT 0,
  total_currency       CHAIN_CURRENCY DEFAULT 0,
  epoch                INTEGER DEFAULT 0,
  slot                 INTEGER DEFAULT 0,
  transactions_count   INTEGER NOT NULL DEFAULT 0,
  transactions_fees    CHAIN_CURRENCY DEFAULT 0,
  snarkers_count       INTEGER NOT NULL DEFAULT 0,
  snark_jobs_count     INTEGER NOT NULL DEFAULT 0,
  snark_jobs_fees      CHAIN_CURRENCY,

  PRIMARY KEY (id)
);

CREATE INDEX idx_blocks_canonic
  ON blocks (canonical);

CREATE INDEX idx_blocks_height
  ON blocks (height);

CREATE INDEX idx_blocks_time
  ON blocks (time);

CREATE UNIQUE INDEX idx_blocks_hash
  ON blocks (hash);

CREATE INDEX idx_blocks_creator
  ON blocks (creator);

-- +goose Down
DROP TABLE IF EXISTS blocks;