CREATE TABLE IF NOT EXISTS delegations (
  id            BIGSERIAL                NOT NULL,
  created_at    TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at    TIMESTAMP WITH TIME ZONE NOT NULL,
  height        DOUBLE PRECISION         NOT NULL,
  time          TIMESTAMP WITH TIME ZONE NOT NULL,
  validator_uid TEXT                     NOT NULL,
  delegator_uid TEXT                     NOT NULL,
  shares        DECIMAL(65, 0)           NOT NULL,

  PRIMARY KEY (time, id)
);

-- Hypertable
SELECT create_hypertable('delegations', 'time', if_not_exists => TRUE);

-- Indexes
CREATE index idx_delegation_height        ON delegations (height, time DESC);
CREATE index idx_delegation_app_version   ON delegations (validator_uid, time DESC);
CREATE index idx_delegation_block_version ON delegations (delegator_uid, time DESC);
