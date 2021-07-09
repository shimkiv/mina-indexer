-- +goose Up
ALTER TABLE blocks ADD COLUMN snarker_accounts TEXT[];

CREATE INDEX idx_snarkers_accounts ON blocks USING GIN (snarker_accounts);

-- +goose Down
ALTER TABLE blocks DROP COLUMN snarker_accounts;

DROP INDEX IF EXISTS idx_snarkers_accounts;