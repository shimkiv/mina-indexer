SELECT
  count(*) as jobs_count,
  sum(works_count) as work_count
FROM
  snark_jobs
WHERE
  prover = $1
  AND block_hash IN
        (SELECT hash FROM blocks
        WHERE height >= $2 AND height <= $3 AND canonical = TRUE
        AND $1 = ANY(snarker_accounts))