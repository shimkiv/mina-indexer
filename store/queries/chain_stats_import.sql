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
  coinbase_amount
)
SELECT
	DATE_TRUNC('@bucket', blocks.time),
	'@bucket',
	ROUND(EXTRACT(EPOCH FROM (MAX(blocks.time) - MIN(blocks.time)) / COUNT(1))::NUMERIC, 2),
	COUNT(1),
	(SELECT COUNT(1) FROM blocks),
	COUNT(DISTINCT(creator)),
	(SELECT COUNT(1) FROM accounts),
	COUNT(DISTINCT(epoch)),
	COUNT(DISTINCT(slot)),
	(SELECT COUNT(1) FROM snarkers),
	SUM(blocks.snark_jobs_count),
  SUM(blocks.snark_jobs_fees),
	AVG(blocks.coinbase),
	AVG(blocks.total_currency),
  COUNT(transactions),
  SUM(transactions.amount),
	COUNT(transactions) FILTER (WHERE type = 'payment'),
	SUM(transactions.amount) FILTER (WHERE type = 'payment'),
  COUNT(transactions) FILTER (WHERE type = 'fee_transfer'),
	SUM(transactions.amount) FILTER (WHERE type = 'fee_transfer'),
  COUNT(transactions) FILTER (WHERE type = 'coinbase'),
	SUM(transactions.amount) FILTER (WHERE type = 'coinbase')
FROM
	blocks
LEFT JOIN transactions
  ON transactions.block_hash = blocks.hash
WHERE
	blocks.time >= $1 AND blocks.time <= $2
GROUP BY
	DATE_TRUNC('@bucket', blocks.time);