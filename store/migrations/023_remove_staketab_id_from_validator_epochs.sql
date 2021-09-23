-- +goose Up
DROP INDEX idx_validator_epochs_account_epoch;

CREATE UNIQUE INDEX idx_validator_epochs_account_epoch
  ON validator_epochs(account_address, epoch);

ALTER TABLE validator_epochs DROP COLUMN staketab_id;

-- +goose Down
ALTER TABLE validator_epochs ADD COLUMN staketab_id INTEGER;

DROP INDEX idx_validator_epochs_account_epoch;

CREATE UNIQUE INDEX idx_validator_epochs_account_epoch
  ON validator_epochs(staketab_id, epoch);

