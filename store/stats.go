package store

import (
	"errors"
	"strings"
	"time"

	"github.com/figment-networks/indexing-engine/store/jsonquery"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/util"
	"github.com/figment-networks/mina-indexer/store/queries"
)

const (
	BucketHour = "h"
	BucketDay  = "d"
)

type StatsStore struct {
	baseStore
}

// CreateChainStats creates a new chain stats record
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
		s.prepareBucket(queries.ChainStatsImport, bucket),
		start, end,
	).Error
}

// CreateValidatorStats creates a new validator stats record
func (s StatsStore) CreateValidatorStats(validator *model.Validator, bucket string, ts time.Time) error {
	start, end, err := s.getTimeRange(bucket, ts)
	if err != nil {
		return err
	}

	return s.db.Exec(
		s.prepareBucket(sqlValidatorStatsImport, bucket),
		start, end, validator.PublicKey,
	).Error
}

// ValidatorStats returns validator stats for a given timeframe
func (s StatsStore) ValidatorStats(validator *model.Validator, period uint, interval string) ([]model.ValidatorStat, error) {
	result := []model.ValidatorStat{}

	err := s.db.
		Model(&model.ValidatorStat{}).
		Where("validator_id = ? AND bucket = ?", validator.ID, interval).
		Order("time DESC").
		Limit(period).
		Find(&result).
		Error

	return result, err
}

// CreateTransactionsStats creates a new transaction stats record
func (s StatsStore) CreateTransactionsStats(bucket string, ts time.Time) error {
	start, end, err := s.getTimeRange(bucket, ts)
	if err != nil {
		return err
	}

	return s.db.Exec(
		s.prepareBucket(queries.TransactionStatsImport, bucket),
		start, end,
	).Error
}

// TransactionsStats returns transactions stats for a given timeframe
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

	sqlValidatorStats = `
		SELECT
			blocks_produced_count,
			delegations_count,
			delegations_amount
		FROM
			validator_stats
		WHERE
			validator_id = $1
			AND bucket = $2
		LIMIT $3`

	sqlValidatorStatsImport = `
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
			(SELECT COUNT(1) FROM accounts WHERE delegate = $3),
			(SELECT COALESCE(SUM(balance::numeric), 0) FROM accounts WHERE delegate = $3)
		)
		ON CONFLICT (time, bucket, validator_id) DO UPDATE
		SET
			blocks_produced_count = excluded.blocks_produced_count,
			delegations_count     = excluded.delegations_count,
			delegations_amount    = excluded.delegations_amount`
)
