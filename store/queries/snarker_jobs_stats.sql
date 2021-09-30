SELECT
  DATE_TRUNC('@bucket', time) AS time,
  '@bucket' AS bucket,
  MIN(fee) fee_min,
  MAX(fee) fee_max,
  ROUND(AVG(fee)) fee_avg,
  COUNT(1) AS jobs_count,
  COUNT(DISTINCT(prover)) AS snarkers_count,
  SUM(works_count) AS works_count
FROM
  snark_jobs
WHERE
  time >= NOW() - INTERVAL '@interval'
  AND time <= NOW()
GROUP BY
  DATE_TRUNC('@bucket', time)
ORDER BY
  time DESC
