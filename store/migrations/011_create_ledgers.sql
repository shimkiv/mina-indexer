-- +goose Up
CREATE TABLE ledgers (
  id                 SERIAL PRIMARY KEY,
  time               CHAIN_TIME,
  epoch              CHAIN_HEIGHT,
  entries_count      INTEGER NOT NULL DEFAULT 0,
  staked_amount      CHAIN_CURRENCY DEFAULT 0,
  delegations_count  INTEGER NOT NULL DEFAULT 0,
  delegations_amount CHAIN_CURRENCY DEFAULT 0
);

CREATE INDEX idx_ledgers_time
  ON ledgers(time);

CREATE UNIQUE INDEX idx_ledgers_epoch
  ON ledgers(epoch);

-- +goose Down
DROP TABLE ledgers;
