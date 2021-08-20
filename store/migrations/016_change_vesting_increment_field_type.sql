-- +goose Up
ALTER TABLE ledger_entries ALTER COLUMN timing_vesting_increment TYPE CHAIN_CURRENCY;

-- +goose Down
ALTER TABLE ledger_entries ALTER COLUMN timing_vesting_increment TYPE INTEGER;