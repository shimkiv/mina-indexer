-- +goose Up
CREATE TABLE IF NOT EXISTS fee_transfers (
  id           SERIAL NOT NULL,
  height       INTEGER NOT NULL,
  time         TIMESTAMP WITH TIME ZONE NOT NULL,
  recipient    TEXT NOT NULL,
  amount       DECIMAL(65, 0) NOT NULL,
  created_at   TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at   TIMESTAMP WITH TIME ZONE NOT NULL,

  PRIMARY KEY (id)
);

CREATE INDEX idx_fee_transfers_height
  ON fee_transfers(height);

CREATE INDEX idx_fee_transfers_time
  ON fee_transfers(time);

CREATE INDEX idx_fee_transfers_recipient
  ON fee_transfers(recipient);

-- +goose Down
DROP TABLE snarkers;
