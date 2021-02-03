-- +goose Up
CREATE TABLE IF NOT EXISTS snarkers (
  id           SERIAL NOT NULL,
  account      TEXT NOT NULL,
  start_height INTEGER NOT NULL,
  start_time   TIMESTAMP WITH TIME ZONE NOT NULL,
  last_height  INTEGER NOT NULL,
  last_time    TIMESTAMP WITH TIME ZONE NOT NULL,
  created_at   TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at   TIMESTAMP WITH TIME ZONE NOT NULL,
  fee          DECIMAL(65, 0) NOT NULL,
  jobs_count   INTEGER DEFAULT 0,
  works_count  INTEGER DEFAULT 0,

  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_snarkers_account
  ON snarkers(account);

-- +goose Down
DROP TABLE snarkers;
