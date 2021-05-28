-- +goose Up
ALTER TABLE ledger_entries ADD COLUMN weight DECIMAL(65, 128);

-- +goose Down
ALTER TABLE ledger_entries DROP COLUMN weight;