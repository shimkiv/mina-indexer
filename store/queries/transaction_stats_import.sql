INSERT INTO transactions_stats (
  time, bucket,
  payments_count, payments_amount,
  delegations_count, delegations_amount,
  block_rewards_count, block_rewards_amount,
  fees_count, fees_amount,
  snark_fees_count, snark_fees_amount
)
SELECT
  DATE_TRUNC('@bucket', time) AS time,
  '@bucket' AS bucket,
  COUNT(1) FILTER (WHERE type = 'payment') AS payments_count,
  COALESCE(SUM(amount) FILTER (WHERE type = 'payment'), 0) AS payments_amount,
  COUNT(1) FILTER (WHERE type = 'delegation') AS delegations_count,
  COALESCE(SUM(amount) FILTER (WHERE type = 'delegation'), 0) AS delegations_amount,
  COUNT(1) FILTER (WHERE type = 'coinbase') AS block_rewards_count,
  COALESCE(SUM(amount) FILTER (WHERE type = 'coinbase'), 0) AS block_rewards_amount,
  COUNT(1) FILTER (WHERE type = 'fee_transfer') AS fees_count,
  COALESCE(SUM(amount) FILTER (WHERE type = 'fee_transfer'), 0) AS fees_amount,
  COUNT(1) FILTER (WHERE type = 'snark_fee') AS snark_fees_count,
  COALESCE(SUM(amount) FILTER (WHERE type = 'snark_fee'), 0) AS snark_fees_amount
FROM
  transactions
WHERE
  time >= $1 AND time <= $2
GROUP BY
  DATE_TRUNC('@bucket', time)
ON CONFLICT (time, bucket) DO UPDATE
SET
  payments_count       = excluded.payments_count,
  payments_amount      = excluded.payments_amount,
  delegations_count    = excluded.delegations_count,
  delegations_amount   = excluded.delegations_amount,
  block_rewards_count  = excluded.block_rewards_count,
  block_rewards_amount = excluded.block_rewards_amount,
  fees_count           = excluded.fees_count,
  fees_amount          = excluded.fees_amount,
  snark_fees_count     = excluded.snark_fees_count,
  snark_fees_amount    = excluded.snark_fees_amount