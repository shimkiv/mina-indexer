CREATE TABLE IF NOT EXISTS transactions (
  id            BIGSERIAL                NOT NULL,
  block_hash    TEXT                     NOT NULL,
  type          TEXT                     NOT NULL,
  hash          TEXT                     NOT NULL,
  height        DOUBLE PRECISION         NOT NULL,
  time          TIMESTAMP WITH TIME ZONE NOT NULL,
  nonce         NUMERIC                  NOT NULL,
  sender_key    TEXT                     NOT NULL,
  recipient_key TEXT                     NOT NULL,
  amount        DECIMAL(65, 0)           NOT NULL,
  fee           DECIMAL(65, 0)           NOT NULL,
  created_at    TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at    TIMESTAMP WITH TIME ZONE NOT NULL,

  PRIMARY KEY (time, id)
);

-- Hypertable
SELECT create_hypertable('transactions', 'time', if_not_exists => TRUE);

-- Indexes
CREATE INDEX idx_transaction_block_hash    ON transactions (block_hash);
CREATE INDEX idx_transaction_height        ON transactions (height, time DESC);
CREATE INDEX idx_transaction_sender_key    ON transactions (sender_key, time DESC);
CREATE INDEX idx_transaction_recipient_key ON transactions (recipient_key, time DESC);
CREATE INDEX idx_transaction_hash          ON transactions (hash, time DESC);