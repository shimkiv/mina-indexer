CREATE TABLE IF NOT EXISTS validators (
  id                  BIGSERIAL                     NOT NULL,
  created_at          TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at          TIMESTAMP WITH TIME ZONE NOT NULL,

  height              DOUBLE PRECISION         NOT NULL,
  time                TIMESTAMP WITH TIME ZONE NOT NULL,
  public_key          TEXT                     NOT NULL,

  PRIMARY KEY (time, id)
);

-- Hypertable
SELECT create_hypertable('validators', 'time', if_not_exists => TRUE);

-- Indexes
CREATE index idx_validators_public_key  ON validators (public_key);
CREATE index idx_validators_height      ON validators (height);
CREATE index idx_validators_height_time ON validators (height, time DESC);
