-- +goose Up
CREATE UNIQUE INDEX idx_ledger_entries_unique
  ON ledger_entries(ledger_id, public_key);

-- +goose Down
DROP INDEX IF EXISTS idx_ledger_entries_unique;
