-- +goose Up
CREATE TABLE transactions_stats (
  id                   SERIAL NOT NULL PRIMARY KEY,
  time                 TIMESTAMP WITH TIME ZONE NOT NULL,
  bucket               e_interval NOT NULL,
  payments_count       INTEGER,
  payments_amount      DECIMAL(65, 0),
  delegations_count    INTEGER,
  delegations_amount   DECIMAL(65, 0),
  fees_count           INTEGER,
  fees_amount          DECIMAL(65, 0),
  snark_fees_count     INTEGER,
  snark_fees_amount    DECIMAL(65, 0),
  block_rewards_count  INTEGER,
  block_rewards_amount DECIMAL(65, 0)
);

CREATE UNIQUE INDEX idx_transactions_stats_bucket
  ON transactions_stats(time, bucket);

-- +goose Down
DROP TABLE transactions_stats;
