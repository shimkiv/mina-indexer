SELECT
  fee,
  COUNT(1) jobs_count,
  SUM(works_count) works_count
FROM snark_jobs
WHERE prover = $1
GROUP BY fee
ORDER BY jobs_count DESC
