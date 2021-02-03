INSERT INTO snarkers (
  account,
  fee,
  jobs_count,
  works_count,
  start_height,
  start_time,
  last_height,
  last_time,
  created_at,
  updated_at
)
VALUES @values
ON CONFLICT (account) DO UPDATE
SET
  fee         = excluded.fee,
  jobs_count  = snarkers.jobs_count + excluded.jobs_count,
  works_count = snarkers.works_count + excluded.works_count,
  last_height = excluded.last_height,
  last_time   = excluded.last_time,
  updated_at  = excluded.updated_at