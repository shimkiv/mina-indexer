-- +goose Up
CREATE TABLE chain_stats (
  time                CHAIN_TIME,
  bucket              CHAIN_INTERVAL,

  block_time_avg      NUMERIC,
  blocks_count        INTEGER,
  blocks_total_count  INTEGER,
  transactions_count  INTEGER,
  validators_count    INTEGER,
  accounts_count      INTEGER,
  epochs_count        INTEGER,
  slots_count         INTEGER,
  snarkers_count      INTEGER,
  jobs_count          INTEGER,
  coinbase            CHAIN_CURRENCY,
  total_currency      CHAIN_CURRENCY,

  PRIMARY KEY (time, bucket)
);

-- +goose Down
DROP TABLE chain_stats;
