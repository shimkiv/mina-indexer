WITH current_delegations AS (
  SELECT * FROM ledger_entries
  WHERE
    ledger_id = (
      SELECT id FROM ledgers
      WHERE epoch = (SELECT MAX(epoch) FROM blocks WHERE time >= $1 AND time <= $2)
      ORDER BY id DESC
      LIMIT 1
    )
    AND delegate = $3
)
INSERT INTO validator_stats (
	time,
	bucket,
	validator_id,
	blocks_produced_count,
	delegations_count,
	delegations_amount
)
VALUES (
	DATE_TRUNC('@bucket', $1::timestamp),
	'@bucket',
	(SELECT id FROM validators WHERE public_key = $3 LIMIT 1),
	(SELECT COUNT(1) FROM blocks WHERE time >= $1 AND time <= $2 AND creator = $3),
  (SELECT COUNT(1) FROM current_delegations),
  (SELECT COALESCE(SUM(balance::numeric), 0) FROM current_delegations)
)
ON CONFLICT (time, bucket, validator_id) DO UPDATE
SET
  blocks_produced_count = excluded.blocks_produced_count,
	delegations_count     = excluded.delegations_count,
	delegations_amount    = excluded.delegations_amount