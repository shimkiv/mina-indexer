-- +goose Up
ALTER TABLE blocks ADD COLUMN supercharged BOOLEAN DEFAULT false;

-- +goose Down
ALTER TABLE blocks DROP COLUMN supercharged;
