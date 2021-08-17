-- +goose Up
ALTER TABLE blocks ADD COLUMN supercharged BOOLEAN;

-- +goose Down
ALTER TABLE blocks DROP COLUMN supercharged;
