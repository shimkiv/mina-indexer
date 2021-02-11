-- +goose Up
CREATE TABLE IF NOT EXISTS accounts (
  id              SERIAL NOT NULL,
  public_key      TEXT NOT NULL,
  delegate        TEXT,
  balance         CHAIN_CURRENCY,
  balance_unknown CHAIN_CURRENCY,
  nonce           INTEGER NOT NULL,
  stake           CHAIN_CURRENCY,
  start_height    CHAIN_HEIGHT,
  start_time      CHAIN_TIME,
  last_height     CHAIN_HEIGHT,
  last_time       CHAIN_TIME,
  created_at      CHAIN_TIME,
  updated_at      CHAIN_TIME,

  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_accounts_public_key
  ON accounts(public_key);

-- +goose Down
DROP TABLE accounts;
