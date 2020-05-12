CREATE TABLE IF NOT EXISTS validators (
  id                  BIGSERIAL                     NOT NULL,
  created_at          TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at          TIMESTAMP WITH TIME ZONE NOT NULL,

  
  height              DOUBLE PRECISION         NOT NULL,
  time                TIMESTAMP WITH TIME ZONE NOT NULL,

  entity_uid          TEXT                     NOT NULL,
  node_uid            TEXT                     NOT NULL,
  consensus_uid       TEXT                     NOT NULL,
  voting_power        DOUBLE PRECISION         NOT NULL,
  total_shares        DECIMAL(65, 0)           NOT NULL,
  proposed            BOOLEAN                  NOT NULL,
  address             TEXT                     NOT NULL,
  precommit_validated BOOLEAN,
  precommit_type      TEXT,
  precommit_index     DOUBLE PRECISION,

  PRIMARY KEY (time, id)
);

-- Hypertable
SELECT create_hypertable('validators', 'time', if_not_exists => TRUE);

-- Indexes
CREATE index idx_validators_height       ON validators (height, time DESC);
CREATE index idx_validators_validator_id ON validators (entity_uid, time DESC);
CREATE index idx_validators_node_uid     ON validators (node_uid, time DESC);
CREATE index idx_validators_proposed     ON validators (proposed, time DESC);
CREATE index idx_validators_total_shares ON validators (total_shares, time DESC);
