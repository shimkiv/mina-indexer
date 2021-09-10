-- +goose Up
CREATE INDEX idx_validator_epochs_account ON validator_epochs(account_address);

-- +goose Down
DROP INDEX idx_validator_epochs_account;
