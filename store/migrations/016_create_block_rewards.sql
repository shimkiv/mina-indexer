-- +goose Up
CREATE TABLE block_rewards (
  id                    CHAIN_UUID,
  owner_account         TEXT NOT NULL,
  delegate              TEXT NOT NULL,
  block_height          CHAIN_HEIGHT NOT NULL,
  block_time            CHAIN_TIME NOT NULL,
  reward                CHAIN_CURRENCY,
  owner_type            OWNER_TYPE
);

CREATE UNIQUE INDEX idx_block_rewards_pk
  ON block_rewards(owner_account, delegate, block_height, owner_type);

CREATE INDEX idx_block_rewards_block_time
  ON block_rewards (block_time);

-- +goose Down
DROP TABLE block_rewards;
