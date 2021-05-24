-- +goose Up
ALTER TABLE blocks ADD COLUMN coinbase_rewards  CHAIN_CURRENCY DEFAULT 0;

-- +goose Down
ALTER TABLE blocks DROP COLUMN coinbase_rewards;