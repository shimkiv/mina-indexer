SELECT
  DATE_TRUNC('d', snark_jobs.time) AS time,
  MIN(fee) AS fee_min,
  MAX(fee) AS fee_max,
  ROUND(AVG(fee)) AS fee_avg,
  COUNT(1) AS jobs_count,
  SUM(works_count) AS works_count,
  COUNT(blocks) AS blocks_count,
  SUM(CASE blocks.canonical WHEN true THEN 1 ELSE 0 END) AS canonical_blocks_count,
  ROUND(SUM(CASE blocks.canonical WHEN true THEN 1 ELSE 0 END) * 100.0 / COUNT(blocks), 2) AS inclusion_ratio,
  SUM(CASE blocks.canonical WHEN true THEN fee ELSE 0 END) AS fees_amount
FROM
  snark_jobs
INNER JOIN blocks
  ON blocks.hash = snark_jobs.block_hash
WHERE
  prover = $1
GROUP BY
  DATE_TRUNC('d', snark_jobs.time)
ORDER BY
  time DESC
LIMIT 30
