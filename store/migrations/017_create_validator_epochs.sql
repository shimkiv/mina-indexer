-- +goose Up
CREATE TABLE validator_epochs (
  id              SERIAL NOT NULL PRIMARY KEY,
  account_id      VARCHAR NOT NULL,
  epoch           INTEGER DEFAULT 0,
  validator_fee   DECIMAL(65, 2) NOT NULL
);

CREATE UNIQUE INDEX idx_validator_epochs_account_epoch
  ON validator_epochs(account_id, epoch);

-- +goose Down
DROP TABLE validator_epochs;
