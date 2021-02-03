INSERT INTO chain_stats (
	time,
	bucket,
	block_time_avg,
	blocks_count,
	blocks_total_count,
	transactions_count,
	validators_count,
	accounts_count,
	epochs_count,
	slots_count,
	snarkers_count,
	jobs_count,
	coinbase,
	total_currency
)
SELECT
	DATE_TRUNC('@bucket', time),
	'@bucket',
	ROUND(EXTRACT(EPOCH FROM (MAX(time) - MIN(time)) / COUNT(1))::NUMERIC, 2),
	COUNT(1),
	(SELECT COUNT(1) FROM blocks),
	SUM(transactions_count),
	COUNT(DISTINCT(creator)),
	(SELECT COUNT(1) FROM accounts),
	COUNT(DISTINCT(epoch)),
	COUNT(DISTINCT(slot)),
	(SELECT COUNT(1) FROM snarkers),
	SUM(snark_jobs_count),
	AVG(coinbase),
	AVG(total_currency)
FROM
	blocks
WHERE
	time >= $1 AND time <= $2
GROUP BY
	DATE_TRUNC('@bucket', time);