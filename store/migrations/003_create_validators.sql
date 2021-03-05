-- +goose Up
CREATE TABLE IF NOT EXISTS validators (
  id                 SERIAL NOT NULL,
  public_key         TEXT NOT NULL,
  identity_name      TEXT,
  start_height       CHAIN_HEIGHT,
  start_time         CHAIN_TIME,
  last_height        CHAIN_HEIGHT,
  last_time          CHAIN_TIME,
  blocks_created     INTEGER NOT NULL DEFAULT 0,
  blocks_proposed    INTEGER NOT NULL DEFAULT 0,
  stake              CHAIN_CURRENCY DEFAULT 0,
  delegations        INTEGER DEFAULT 0,
  created_at         CHAIN_TIME,
  updated_at         CHAIN_TIME,

  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_validators_publix_key
  ON validators (public_key);

-- +goose Down
DROP TABLE validators;
