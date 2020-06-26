package store

import (
	"errors"
	"strings"
	"time"

	"github.com/figment-networks/coda-indexer/model/util"
	"github.com/figment-networks/indexing-engine/store/jsonquery"
)

const (
	BucketHour = "h"
	BucketDay  = "d"
)

type StatsStore struct {
	baseStore
}

func (s StatsStore) CreateChainStats(bucket string, ts time.Time) error {
	start, end, err := s.getTimeRange(bucket, ts)
	if err != nil {
		return err
	}

	err = s.db.Exec(
		s.prepareBucket(sqlChainStatsDelete, bucket),
		start,
	).Error
	if err != nil && err != ErrNotFound {
		return err
	}

	return s.db.Exec(
		s.prepareBucket(sqlChainStatsImport, bucket),
		start, end,
	).Error
}

func (s StatsStore) CreateTransactionsStats(bucket string, ts time.Time) error {
	start, end, err := s.getTimeRange(bucket, ts)
	if err != nil {
		return err
	}

	return s.db.Exec(
		s.prepareBucket(sqlTransactionsStatsImport, bucket),
		start, end,
	).Error
}

func (s StatsStore) TransactionsStats(period uint, interval string) ([]byte, error) {
	return jsonquery.MustArray(s.db, sqlTransactionsStats, period, interval)
}

// getTimeRange returns the start/end time for a given time bucket
func (s StatsStore) getTimeRange(bucket string, ts time.Time) (start time.Time, end time.Time, err error) {
	switch bucket {
	case BucketHour:
		start, end = util.HourInterval(ts)
	case BucketDay:
		start, end = util.DayInterval(ts)
	default:
		err = errors.New("invalid time bucket")
	}
	return
}

func (s StatsStore) prepareBucket(q, bucket string) string {
	return strings.ReplaceAll(q, "@bucket", bucket)
}

var (
	sqlChainStatsDelete = `DELETE FROM chain_stats WHERE time = ? AND BUCKET = '@bucket';`

	sqlChainStatsImport = `
		INSERT INTO chain_stats (
			time,
			bucket,
			block_time_avg,
			blocks_count,
			transactions_count,
			fee_transfers_count,
			validators_count,
			accounts_count,
			epochs_count,
			slots_count,
			snarkers_count, snarkers_avg, snarkers_min, snarkers_max,
			jobs_count, jobs_min, jobs_max, jobs_avg,
			coinbase_max, coinbase_min, coinbase_diff,
			total_currency_max, total_currency_min, total_currency_diff
		)
		SELECT
			DATE_TRUNC('@bucket', time) AS time,
			'@bucket' AS bucket,

			ROUND(EXTRACT(EPOCH FROM (MAX(time) - MIN(time)) / COUNT(1))::NUMERIC, 2) AS block_time_avg,
			COUNT(1) AS blocks_count,

			SUM(transactions_count) AS transactions_count,
			SUM(fee_transfers_count) AS fee_transfers_count,

			COUNT(DISTINCT(creator)) AS validators_count,
			(SELECT COUNT(1) FROM accounts) AS accounts_count,
			COUNT(DISTINCT(epoch)) AS epochs_count,
			COUNT(DISTINCT(slot)) AS slots_count,

			(SELECT COUNT(1) FROM snarkers) AS snarkers_count,
			ROUND(AVG(snarkers_count), 4) AS snarkers_avg,
			MIN(snarkers_count) AS snarkers_min,
			MAX(snarkers_count) AS snarkers_max,

			SUM(snark_jobs_count) AS jobs_count,
			MIN(snark_jobs_count) AS jobs_min,
			MAX(snark_jobs_count) AS jobs_max,
			ROUND(AVG(snark_jobs_count), 4) AS jobs_avg,

			MAX(coinbase) AS coinbase_max,
			MIN(coinbase) AS coinbase_min,
			(MAX(coinbase) - MIN(coinbase)) AS coinbase_diff,

			MAX(total_currency) AS total_currency_max,
			MIN(total_currency) AS total_currency_min,
			(MAX(total_currency) - MIN(total_currency)) AS total_currency_diff
		FROM
			blocks
		WHERE
			time >= $1 AND time <= $2
		GROUP BY
			DATE_TRUNC('@bucket', time);`

	sqlTransactionsStatsImport = `
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
  		COUNT(1) FILTER (WHERE type = 'block_reward') AS block_rewards_count,
  		COALESCE(SUM(amount) FILTER (WHERE type = 'block_reward'), 0) AS block_rewards_amount,
  		COUNT(1) FILTER (WHERE type = 'fee') AS fees_count,
  		COALESCE(SUM(amount) FILTER (WHERE type = 'fee'), 0) AS fees_amount,
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
			snark_fees_amount    = excluded.snark_fees_amount`

	sqlTransactionsStats = `
		SELECT
			time,
			payments_count,
			payments_amount,
			delegations_count,
			delegations_amount,
			block_rewards_count,
			block_rewards_amount,
			fees_count,
			fees_amount,
			snark_fees_count,
			snark_fees_amount
		FROM
			transactions_stats
		WHERE
			bucket = $2
		ORDER BY
			time DESC
		LIMIT $1`
)
