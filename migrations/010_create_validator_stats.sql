-- +goose Up
CREATE TABLE validator_stats (
  id                    SERIAL NOT NULL PRIMARY KEY,
  time                  TIMESTAMP WITH TIME ZONE NOT NULL,
  bucket                e_interval NOT NULL,
  validator_id          INTEGER,
  blocks_produced_count INTEGER,
  delegations_count     INTEGER,
  delegations_amount    DECIMAL(65, 0)
);

CREATE UNIQUE INDEX idx_validator_stats_bucket
  ON validator_stats(time, bucket, validator_id);

-- +goose Down
DROP TABLE validator_stats;
