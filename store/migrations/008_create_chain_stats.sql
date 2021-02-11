-- +goose Up
CREATE TABLE chain_stats (
  time                 CHAIN_TIME,
  bucket               CHAIN_INTERVAL,
  block_time_avg       NUMERIC DEFAULT 0,
  blocks_count         INTEGER DEFAULT 0,
  blocks_total_count   INTEGER DEFAULT 0,
  validators_count     INTEGER DEFAULT 0,
  accounts_count       INTEGER DEFAULT 0,
  epochs_count         INTEGER DEFAULT 0,
  slots_count          INTEGER DEFAULT 0,
  snarkers_count       INTEGER DEFAULT 0,
  jobs_count           INTEGER DEFAULT 0,
  jobs_amount          CHAIN_CURRENCY DEFAULT 0,
  coinbase             CHAIN_CURRENCY DEFAULT 0,
  total_currency       CHAIN_CURRENCY DEFAULT 0,
  transactions_count   INTEGER DEFAULT 0,
  transactions_amount  CHAIN_CURRENCY DEFAULT 0,
  payments_count       INTEGER DEFAULT 0,
  payments_amount      CHAIN_CURRENCY DEFAULT 0,
  fee_transfers_count  INTEGER DEFAULT 0,
  fee_transfers_amount CHAIN_CURRENCY DEFAULT 0,
  coinbase_count       INTEGER DEFAULT 0,
  coinbase_amount      CHAIN_CURRENCY DEFAULT 0,

  PRIMARY KEY (time, bucket)
);

-- +goose Down
DROP TABLE chain_stats;
