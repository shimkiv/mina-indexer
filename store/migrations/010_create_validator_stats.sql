-- +goose Up
CREATE TABLE validator_stats (
  id                    CHAIN_UUID,
  validator_id          INTEGER NOT NULL,
  time                  CHAIN_TIME,
  bucket                CHAIN_INTERVAL,
  blocks_produced_count INTEGER DEFAULT 0,
  delegations_count     INTEGER DEFAULT 0,
  delegations_amount    CHAIN_CURRENCY
);

CREATE UNIQUE INDEX idx_validator_stats_bucket
  ON validator_stats(time, bucket, validator_id);

-- +goose Down
DROP TABLE validator_stats;
