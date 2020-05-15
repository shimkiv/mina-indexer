CREATE TABLE IF NOT EXISTS accounts (
  id           BIGSERIAL                NOT NULL,
  start_height NUMERIC                  NOT NULL,
  public_key   TEXT                     NOT NULL,
  balance      DECIMAL(65, 0)           NOT NULL,
  nonce        NUMERIC                  NOT NULL,
  started_at   TIMESTAMP WITH TIME ZONE NOT NULL,
  created_at   TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at   TIMESTAMP WITH TIME ZONE NOT NULL,

  PRIMARY KEY (id)
);

-- Indexes
CREATE INDEX idx_accounts_public_key ON accounts (public_key);
CREATE INDEX idx_accounts_height     ON accounts (start_height);