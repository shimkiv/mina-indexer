package queries

const (

	// AccountsImport is imported from store/queries/accounts_import.sql
	AccountsImport = `INSERT INTO accounts (   public_key,   delegate,   balance,   balance_unknown,   nonce,   start_height,   start_time,   last_height,   last_time,   created_at,   updated_at ) VALUES @values ON CONFLICT (public_key) DO UPDATE SET   delegate        = excluded.delegate,   balance         = excluded.balance,   balance_unknown = excluded.balance_unknown,   nonce           = excluded.nonce,   last_height     = excluded.last_height,   last_time       = excluded.last_time,   updated_at      = excluded.updated_at`

	// BlocksStats is imported from store/queries/blocks_stats.sql
	BlocksStats = `SELECT   time,   block_time_avg,   blocks_count,   validators_count,   snarkers_count,   transactions_count,   jobs_count FROM   chain_stats WHERE   bucket = $2 ORDER BY   time DESC LIMIT $1`

	// BlocksTimes is imported from store/queries/blocks_times.sql
	BlocksTimes = `SELECT   MIN(height) start_height,   MAX(height) end_height,   MIN(time) start_time,   MAX(time) end_time,   COUNT(*) count,   EXTRACT(EPOCH FROM MAX(time) - MIN(time)) AS diff,   EXTRACT(EPOCH FROM ((MAX(time) - MIN(time)) / COUNT(*))) AS avg FROM   (     SELECT * FROM blocks     ORDER BY height DESC     LIMIT ?   ) t`

	// ChainStatsImport is imported from store/queries/chain_stats_import.sql
	ChainStatsImport = `INSERT INTO chain_stats ( 	time, 	bucket, 	block_time_avg, 	blocks_count, 	blocks_total_count, 	transactions_count, 	validators_count, 	accounts_count, 	epochs_count, 	slots_count, 	snarkers_count, 	jobs_count, 	coinbase, 	total_currency ) SELECT 	DATE_TRUNC('@bucket', time), 	'@bucket', 	ROUND(EXTRACT(EPOCH FROM (MAX(time) - MIN(time)) / COUNT(1))::NUMERIC, 2), 	COUNT(1), 	(SELECT COUNT(1) FROM blocks), 	SUM(transactions_count), 	COUNT(DISTINCT(creator)), 	(SELECT COUNT(1) FROM accounts), 	COUNT(DISTINCT(epoch)), 	COUNT(DISTINCT(slot)), 	(SELECT COUNT(1) FROM snarkers), 	SUM(snark_jobs_count), 	AVG(coinbase), 	AVG(total_currency) FROM 	blocks WHERE 	time >= $1 AND time <= $2 GROUP BY 	DATE_TRUNC('@bucket', time);`

	// SnarkJobsImport is imported from store/queries/snark_jobs_import.sql
	SnarkJobsImport = `INSERT INTO snark_jobs (   height,   time,   prover,   fee,   works_count,   created_at ) VALUES @values`

	// SnarkersImport is imported from store/queries/snarkers_import.sql
	SnarkersImport = `INSERT INTO snarkers (   account,   fee,   jobs_count,   works_count,   start_height,   start_time,   last_height,   last_time,   created_at,   updated_at ) VALUES @values ON CONFLICT (account) DO UPDATE SET   fee         = excluded.fee,   jobs_count  = snarkers.jobs_count + excluded.jobs_count,   works_count = snarkers.works_count + excluded.works_count,   last_height = excluded.last_height,   last_time   = excluded.last_time,   updated_at  = excluded.updated_at`

	// TransactionStatsImport is imported from store/queries/transaction_stats_import.sql
	TransactionStatsImport = `INSERT INTO transactions_stats (   time, bucket,   payments_count, payments_amount,   delegations_count, delegations_amount,   block_rewards_count, block_rewards_amount,   fees_count, fees_amount,   snark_fees_count, snark_fees_amount ) SELECT   DATE_TRUNC('@bucket', time) AS time,   '@bucket' AS bucket,   COUNT(1) FILTER (WHERE type = 'payment') AS payments_count,   COALESCE(SUM(amount) FILTER (WHERE type = 'payment'), 0) AS payments_amount,   COUNT(1) FILTER (WHERE type = 'delegation') AS delegations_count,   COALESCE(SUM(amount) FILTER (WHERE type = 'delegation'), 0) AS delegations_amount,   COUNT(1) FILTER (WHERE type = 'coinbase') AS block_rewards_count,   COALESCE(SUM(amount) FILTER (WHERE type = 'coinbase'), 0) AS block_rewards_amount,   COUNT(1) FILTER (WHERE type = 'fee_transfer') AS fees_count,   COALESCE(SUM(amount) FILTER (WHERE type = 'fee_transfer'), 0) AS fees_amount,   COUNT(1) FILTER (WHERE type = 'snark_fee') AS snark_fees_count,   COALESCE(SUM(amount) FILTER (WHERE type = 'snark_fee'), 0) AS snark_fees_amount FROM   transactions WHERE   time >= $1 AND time <= $2 GROUP BY   DATE_TRUNC('@bucket', time) ON CONFLICT (time, bucket) DO UPDATE SET   payments_count       = excluded.payments_count,   payments_amount      = excluded.payments_amount,   delegations_count    = excluded.delegations_count,   delegations_amount   = excluded.delegations_amount,   block_rewards_count  = excluded.block_rewards_count,   block_rewards_amount = excluded.block_rewards_amount,   fees_count           = excluded.fees_count,   fees_amount          = excluded.fees_amount,   snark_fees_count     = excluded.snark_fees_count,   snark_fees_amount    = excluded.snark_fees_amount`

	// TransactionsImport is imported from store/queries/transactions_import.sql
	TransactionsImport = `INSERT INTO transactions (   type,   hash,   block_hash,   block_height,   time,   nonce,   sender,   receiver,   amount,   fee,   memo,   status,   failure_reason,   sequence_number,   secondary_sequence_number,   created_at,   updated_at ) VALUES @values  ON CONFLICT (hash) DO UPDATE SET   updated_at = excluded.updated_at`

	// ValidatorsImport is imported from store/queries/validators_import.sql
	ValidatorsImport = `INSERT INTO validators (   public_key,   start_height,   start_time,   last_height,   last_time,   stake,   blocks_proposed,   blocks_created,   delegated_accounts,   delegated_balance,   created_at,   updated_at ) VALUES   @values ON CONFLICT (public_key) DO UPDATE SET   last_height        = excluded.last_height,   last_time          = excluded.last_time,   stake              = excluded.stake,   blocks_proposed    = excluded.blocks_proposed,   blocks_created     = COALESCE((SELECT COUNT(1) FROM blocks WHERE creator = excluded.public_key), 0),   delegated_accounts = COALESCE((SELECT COUNT(1) FROM accounts WHERE delegate = excluded.public_key), 0),   delegated_balance  = COALESCE((SELECT SUM(balance::NUMERIC) FROM accounts WHERE delegate = excluded.public_key), 0),   updated_at         = excluded.updated_at`

	// ValidatorsIndex is imported from store/queries/validators_index.sql
	ValidatorsIndex = `SELECT   validators.public_key,   validators.identity_name,   validators.start_height,   validators.start_time,   validators.last_height,   validators.last_time,   COALESCE(validators.stake::TEXT) AS stake,   validators.blocks_created,   validators.blocks_proposed,   validators.delegated_accounts,   COALESCE(validators.delegated_balance, 0)::TEXT AS delegated_balance,   COALESCE(accounts.balance, 0)::TEXT AS account_balance,   COALESCE(accounts.balance_unknown, 0)::TEXT AS account_balance_unknown FROM   validators LEFT JOIN accounts   ON accounts.public_key = validators.public_key ORDER BY   blocks_created DESC`
)
