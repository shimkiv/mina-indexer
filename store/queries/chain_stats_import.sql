WITH current_ledger AS (
  SELECT * FROM ledger_entries
  WHERE ledger_id = (
    SELECT id FROM ledgers
    WHERE epoch = (SELECT MAX(epoch) FROM blocks WHERE time >= $1 AND time <= $2)
    LIMIT 1
  )
)
INSERT INTO chain_stats (
  time,
  bucket,
  block_time_avg,
  blocks_count,
  blocks_total_count,
  validators_count,
  accounts_count,
  epochs_count,
  slots_count,
  snarkers_count,
  jobs_count,
  jobs_amount,
  coinbase,
  total_currency,
  transactions_count,
  transactions_amount,
  payments_count,
  payments_amount,
  fee_transfers_count,
  fee_transfers_amount,
  coinbase_count,
  coinbase_amount,
  staked_amount,
  delegations_count,
  delegations_amount
)
SELECT
  DATE_TRUNC('@bucket', blocks.time),
  '@bucket',
  ROUND(EXTRACT(EPOCH FROM (MAX(blocks.time) - MIN(blocks.time)) / COUNT(DISTINCT blocks.height))::NUMERIC, 2),
  COUNT(DISTINCT blocks.height),
  (SELECT COUNT(1) FROM blocks),
  COUNT(DISTINCT(creator)),
  (SELECT COUNT(1) FROM accounts),
  COUNT(DISTINCT(blocks.epoch)),
  COUNT(DISTINCT(blocks.slot)),
  (SELECT COUNT(1) FROM snarkers),
  COALESCE(SUM(blocks.snark_jobs_count), 0),
  COALESCE(SUM(blocks.snark_jobs_fees), 0),
  COALESCE(AVG(blocks.coinbase), 0),
  COALESCE(AVG(blocks.total_currency), 0),
  COUNT(transactions),
  COALESCE(SUM(transactions.amount), 0),
  COUNT(transactions) FILTER (WHERE type = 'payment'),
  COALESCE(SUM(transactions.amount) FILTER (WHERE type = 'payment'), 0),
  COUNT(transactions) FILTER (WHERE type = 'fee_transfer'),
  COALESCE(SUM(transactions.amount) FILTER (WHERE type = 'fee_transfer'), 0),
  COUNT(transactions) FILTER (WHERE type = 'coinbase'),
  COALESCE(SUM(transactions.amount) FILTER (WHERE type = 'coinbase'), 0),
  COALESCE((SELECT SUM(balance) FROM current_ledger), 0),
  COALESCE((SELECT COUNT(1) FROM current_ledger WHERE delegation IS TRUE), 0),
  COALESCE((SELECT SUM(balance) FROM current_ledger WHERE delegation IS TRUE), 0)
FROM
  blocks
LEFT JOIN transactions
  ON transactions.block_hash = blocks.hash
  AND transactions.status = 'applied'
WHERE
  blocks.time >= $1
  AND blocks.time <= $2
  AND blocks.canonical = true
  AND transactions.canonical = true
GROUP BY
  DATE_TRUNC('@bucket', blocks.time);