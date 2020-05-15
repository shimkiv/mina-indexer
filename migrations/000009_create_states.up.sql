CREATE TABLE IF NOT EXISTS states (
  id                                    BIGSERIAL                NOT NULL,
  created_at                            TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at                            TIMESTAMP WITH TIME ZONE NOT NULL,
  height                                DOUBLE PRECISION         NOT NULL,
  total_currency                        DECIMAL(65, 0)           NOT NULL,
  epoch                                 DOUBLE PRECISION,
  epoch_count                           DOUBLE PRECISION,
  last_vfr_output                       TEXT,

  PRIMARY KEY (id)
);

-- Indexes
CREATE INDEX idx_states_height ON states (height);