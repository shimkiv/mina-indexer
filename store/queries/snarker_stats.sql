SELECT
  ROUND(SUM(CASE blocks.canonical WHEN true THEN 1 ELSE 0 END) * 100.0 / COUNT(blocks), 2) AS inclusion_ratio,
  ROUND(COUNT(1) * 100.0 / (SELECT COUNT(1) FROM snark_jobs), 4) AS jobs_percent,
  ROUND(SUM(CASE blocks.canonical WHEN true THEN fee ELSE 0 END)) AS fees_amount
FROM
  snark_jobs
INNER JOIN blocks
  ON blocks.hash = snark_jobs.block_hash
WHERE
  prover = ?
