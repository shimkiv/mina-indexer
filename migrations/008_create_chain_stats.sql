-- +goose Up
CREATE TYPE e_interval AS ENUM ('h', 'd', 'w');

CREATE TABLE chain_stats (
  id                  SERIAL NOT NULL PRIMARY KEY,
  time                TIMESTAMP WITH TIME ZONE NOT NULL,
  bucket              e_interval NOT NULL,

  block_time_avg      NUMERIC,
  blocks_count        INTEGER,
  blocks_total_count  INTEGER,
  transactions_count  INTEGER,
  fee_transfers_count INTEGER,
  validators_count    INTEGER,
  accounts_count      INTEGER,
  epochs_count        INTEGER,
  slots_count         INTEGER,

  snarkers_count INTEGER,
  snarkers_avg   NUMERIC,
  snarkers_min   INTEGER,
  snarkers_max   INTEGER,

  jobs_count INTEGER,
  jobs_min   INTEGER,
  jobs_max   INTEGER,
  jobs_avg   INTEGER,

  coinbase_max  DECIMAL(65, 0),
  coinbase_min  DECIMAL(65, 0),
  coinbase_diff DECIMAL(65, 0),

  total_currency_max  DECIMAL(65, 0),
  total_currency_min  DECIMAL(64, 0),
  total_currency_diff DECIMAL(65, 0)
);

CREATE UNIQUE INDEX idx_chain_stats_bucket
  ON chain_stats(time, bucket);

-- +goose Down
DROP TABLE chain_stats;
DROP TYPE e_interval;
