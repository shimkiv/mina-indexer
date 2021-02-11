-- +goose Up
CREATE TYPE CHAIN_TX_TYPE AS ENUM (
  'payment',
  'delegation',
  'coinbase',
  'fee_transfer_via_coinbase',
  'fee_transfer',
  'snark_fee'
);

CREATE TYPE CHAIN_TX_STATUS AS ENUM ('applied', 'failed');

CREATE TABLE transactions (
  id                        SERIAL NOT NULL,
  type                      CHAIN_TX_TYPE,
  hash                      TEXT NOT NULL,
  block_hash                TEXT NOT NULL,
  block_height              CHAIN_HEIGHT,
  time                      CHAIN_TIME,
  nonce                     INTEGER,
  sender                    TEXT,
  receiver                  TEXT NOT NULL,
  amount                    CHAIN_CURRENCY,
  fee                       CHAIN_CURRENCY,
  memo                      TEXT,
  status                    CHAIN_TX_STATUS,
  failure_reason            TEXT,
  sequence_number           INTEGER,
  secondary_sequence_number INTEGER,
  created_at                CHAIN_TIME,
  updated_at                CHAIN_TIME
);

CREATE INDEX idx_transactions_type
  ON transactions(type);

CREATE INDEX idx_transactions_block_hash
  ON transactions(block_hash);

CREATE UNIQUE INDEX idx_transactions_hash
  ON transactions(hash);

CREATE INDEX idx_transactions_block_height
  ON transactions(block_height);

CREATE INDEX idx_transactions_time
  ON transactions(time);

CREATE INDEX idx_transactions_sender
  ON transactions(sender);

CREATE INDEX idx_transactions_receiver
  ON transactions(receiver);

CREATE INDEX idx_transactions_memo
  ON transactions(LOWER(memo));

-- +goose Down
DROP TABLE transactions;
DROP TYPE CHAIN_TX_TYPE;
DROP TYPE CHAIN_TX_STATUS;
