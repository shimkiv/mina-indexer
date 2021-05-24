-- +goose Up
ALTER TABLE ledger_entries ADD COLUMN weight NUMERIC;

-- +goose Down
ALTER TABLE ledger_entries DROP COLUMN weight;