-- +goose Up
CREATE TABLE block_rewards (
  id                    CHAIN_UUID,
  owner_account         TEXT NOT NULL,
  delegate              TEXT NOT NULL,
  time_bucket           CHAIN_TIME,
  reward                NUMERIC(65, 4) DEFAULT 0,
  owner_type            OWNER_TYPE
);

CREATE UNIQUE INDEX idx_block_rewards_pk
  ON block_rewards(owner_account, delegate, time_bucket, owner_type);

CREATE INDEX idx_block_rewards_time_bucket
  ON block_rewards (time_bucket);

-- +goose Down
DROP TABLE block_rewards;
