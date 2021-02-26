SELECT
  time,
  block_time_avg,
  blocks_count,
  validators_count,
  snarkers_count,
  jobs_count,
  jobs_amount::TEXT jobs_amount,
  transactions_count,
  transactions_amount::TEXT transactions_amount,
  payments_count,
  payments_amount::TEXT payments_amount,
  fee_transfers_count,
  fee_transfers_amount::TEXT fee_transfers_amount,
  coinbase_count,
  coinbase_amount::TEXT coinbase_amount,
  total_currency::TEXT total_currency
FROM
  chain_stats
WHERE
  bucket = $2
ORDER BY
  time DESC
LIMIT
  $1