-- +goose Up
ALTER TABLE blocks ADD COLUMN reward_calculated BOOLEAN;

-- +goose Down
ALTER TABLE blocks DROP COLUMN reward_calculated;