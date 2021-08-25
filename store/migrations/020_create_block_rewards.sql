-- +goose Up
CREATE TABLE block_rewards (
  id                    CHAIN_UUID,
  owner_account         TEXT NOT NULL,
  delegate              TEXT,
  epoch                 INTEGER NOT NULL,
  time_bucket           CHAIN_TIME,
  reward                CHAIN_CURRENCY,
  owner_type            OWNER_TYPE
);

CREATE UNIQUE INDEX idx_block_rewards_pk
  ON block_rewards(owner_account, delegate, epoch, time_bucket, owner_type);

-- +goose Down
DROP TABLE block_rewards;
