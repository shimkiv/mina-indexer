SELECT
  time,
  block_time_avg,
  blocks_count,
  validators_count,
  snarkers_count,
  transactions_count,
  jobs_count
FROM
  chain_stats
WHERE
  bucket = $2
ORDER BY
  time DESC
LIMIT $1