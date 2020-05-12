CREATE TABLE IF NOT EXISTS blocks (
  id                  BIGSERIAL                NOT NULL,
  created_at          TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at          TIMESTAMP WITH TIME ZONE NOT NULL,
  height              DOUBLE PRECISION         NOT NULL,
  time                TIMESTAMP WITH TIME ZONE NOT NULL,
  hash                TEXT                     NOT NULL,
  parent_hash         TEXT                     NOT NULL,
  ledger_hash         TEXT                     NOT NULL,
  creator             TEXT                     NOT NULL,
  transactions_count  DOUBLE PRECISION         NOT NULL,
  coinbase            DECIMAL(65, 0)           NOT NULL,
  app_version         TEXT                     NOT NULL,

  PRIMARY KEY (time, id)
);

-- Hypertable
SELECT create_hypertable('blocks', 'time', if_not_exists => TRUE);

-- Indexes
CREATE index idx_blocks_height      ON blocks (height, time DESC);
CREATE index idx_blocks_app_version ON blocks (app_version, time DESC);
CREATE index idx_blocks_hash        ON blocks (hash, time DESC);
CREATE index idx_blocks_creator     ON blocks (creator, time DESC);
