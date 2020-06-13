package store

import (
	"errors"
	"strings"
	"time"

	"github.com/figment-networks/coda-indexer/model/util"
)

const (
	BucketHour = "h"
	BucketDay  = "d"
)

type StatsStore struct {
	baseStore
}

func (s StatsStore) CreateChainStats(bucket string, ts time.Time) error {
	var start, end time.Time

	switch bucket {
	case BucketHour:
		start, end = util.HourInterval(ts)
	case BucketDay:
		start, end = util.DayInterval(ts)
	default:
		return errors.New("invalid time bucket")
	}

	err := s.db.Exec(
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

			ROUND(EXTRACT(EPOCH FROM (LAST(time, time) - FIRST(time, time)) / COUNT(1))::NUMERIC, 2) AS block_time_avg,
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
)
