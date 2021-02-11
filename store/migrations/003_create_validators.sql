-- +goose Up
CREATE TABLE IF NOT EXISTS validators (
  id                 SERIAL NOT NULL,
  public_key         TEXT NOT NULL,
  identity_name      TEXT,
  start_height       CHAIN_HEIGHT,
  start_time         CHAIN_TIME,
  last_height        CHAIN_HEIGHT,
  last_time          CHAIN_TIME,
  stake              CHAIN_CURRENCY,
  blocks_created     INTEGER NOT NULL DEFAULT 0,
  blocks_proposed    INTEGER NOT NULL DEFAULT 0,
  delegated_accounts INTEGER NOT NULL DEFAULT 0,
  delegated_balance  CHAIN_CURRENCY,
  created_at         CHAIN_TIME,
  updated_at         CHAIN_TIME,

  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_validators_publix_key
  ON validators (public_key);

-- +goose Down
DROP TABLE validators;
