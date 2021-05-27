-- +goose Up
CREATE TABLE block_rewards (
  id                    CHAIN_UUID,
  public_key            TEXT NOT NULL,
  delegate              TEXT NOT NULL,
  block_height          CHAIN_HEIGHT NOT NULL,
  block_time            CHAIN_TIME NOT NULL,
  reward                CHAIN_CURRENCY,
  reward_owner_type     REWARD_OWNER_TYPE
);

CREATE UNIQUE INDEX idx_block_rewards_pk
  ON block_rewards(public_key, delegate, block_height, reward_owner_type);

CREATE INDEX idx_block_rewards_block_time
  ON block_rewards (block_time);

-- +goose Down
DROP TABLE block_rewards;
