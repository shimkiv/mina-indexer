-- +goose Up
CREATE TABLE IF NOT EXISTS validators (
  id                 SERIAL NOT NULL,
  account            TEXT NOT NULL,
  start_height       INTEGER NOT NULL,
  start_time         TIMESTAMP WITH TIME ZONE NOT NULL,
  last_height        INTEGER NOT NULL,
  last_time          TIMESTAMP WITH TIME ZONE NOT NULL,
  blocks_created     INTEGER DEFAULT 0,
  blocks_proposed    INTEGER DEFAULT 0,
  delegated_accounts INTEGER DEFAULT 0,
  delegated_balance  DECIMAL(24, 0),
  created_at         TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at         TIMESTAMP WITH TIME ZONE NOT NULL,

  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_validators_account
  ON validators (account);

-- +goose Down
DROP TABLE validators;
