-- +goose Up
CREATE TABLE validator_epochs
(
    id              SERIAL PRIMARY KEY,
    staketab_id     INTEGER        NOT NULL,
    account_address VARCHAR        NOT NULL,
    epoch           INTEGER        NOT NULL,
    validator_fee   DECIMAL(65, 2) NOT NULL
);

CREATE UNIQUE INDEX idx_validator_epochs_account_epoch
  ON validator_epochs(staketab_id, epoch);

-- +goose Down
DROP TABLE validator_epochs;
