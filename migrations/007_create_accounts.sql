-- +goose Up
CREATE TABLE IF NOT EXISTS accounts (
  id              SERIAL NOT NULL,
  public_key      TEXT NOT NULL,
  name            TEXT,
  delegate        TEXT,
  balance         TEXT NOT NULL,
  balance_unknown TEXT NOT NULL,
  nonce           INTEGER NOT NULL,
  start_height    INTEGER NOT NULL,
  start_time      TIMESTAMP WITH TIME ZONE NOT NULL,
  last_height     INTEGER NOT NULL,
  last_time       TIMESTAMP WITH TIME ZONE NOT NULL,
  created_at      TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at      TIMESTAMP WITH TIME ZONE NOT NULL,

  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_accounts_public_key
  ON accounts(public_key);

-- +goose Down
DROP TABLE accounts;
