package store

import (
	"errors"
	"strings"
	"time"

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
func (s StatsStore) CreateValidatorStats(validatorPublicKey string, bucket string, ts time.Time) error {
	start, end, err := s.getTimeRange(bucket, ts)
	if err != nil {
		return err
	}

	return s.db.Exec(
		s.prepareBucket(queries.ValidatorsCreateStats, bucket),
		start, end, validatorPublicKey,
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

// FindValidatorsForDefaultStats returns validator for default values
func (s StatsStore) FindValidatorsForDefaultStats(bucket string, ts time.Time) ([]model.Validator, error) {
	start, _, err := s.getTimeRange(bucket, ts)
	if err != nil {
		return nil, err
	}

	var res []model.Validator
	err = s.db.Raw(queries.ValidatorsForDefaultStats, bucket, start).Scan(&res).Error
	if err != nil {
		return nil, checkErr(err)
	}
	return res, nil
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
)
