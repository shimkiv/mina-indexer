-- +goose Up
CREATE TABLE delegator_block_rewards (
  id                    CHAIN_UUID,
  public_key            TEXT NOT NULL,
  block_height          CHAIN_HEIGHT,
  block_time            CHAIN_TIME,
  delegate              TEXT,
  reward                CHAIN_CURRENCY
);

CREATE UNIQUE INDEX idx_delegator_block_rewards_pk
  ON delegator_block_rewards(public_key, delegate, block_time);

CREATE INDEX idx_delegator_block_rewards_block_time ON delegator_block_rewards (block_time);

-- +goose Down
DROP TABLE delegator_block_rewards;
