package store

import (
	"strings"
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"

	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/store/queries"
)

// RewardStore handles operations on rewards
type RewardStore struct {
	baseStore
}

// FetchRewardsByInterval fetches rewards by interval
func (s *RewardStore) FetchRewardsByInterval(
	ownerAccount string,
	delegate string,
	from time.Time,
	to time.Time,
	timeInterval model.TimeInterval,
	rewardOwnerType model.RewardOwnerType,
	includeEpoch bool,
) ([]model.RewardsSummary, error) {
	var (
		sqlSelect string
		sqlGroup  string
	)

	if includeEpoch {
		sqlSelect = "to_char(time_bucket, $INTERVAL) AS interval, delegate,  epoch,  SUM(reward) AS amount"
	} else {
		sqlSelect = "to_char(time_bucket, $INTERVAL) AS interval, delegate, SUM(reward) AS amount"
	}
	sqlSelect = strings.Replace(sqlSelect, "$INTERVAL", "'"+timeInterval.String()+"'", -1)

	scope := s.db.
		Table("block_rewards").
		Select(sqlSelect)

	if ownerAccount != "" {
		scope = scope.Where("owner_account = ?", ownerAccount)
	}
	if delegate != "" {
		scope = scope.Where("delegate = ?", delegate)
	}
	if rewardOwnerType != "" {
		scope = scope.Where("owner_type = ?", rewardOwnerType)
	}
	if !from.IsZero() {
		scope = scope.Where("time_bucket >= ?", from)
	}
	if !to.IsZero() {
		scope = scope.Where("time_bucket < ?", to)
	}

	if includeEpoch {
		sqlGroup = "to_char(time_bucket, $INTERVAL), delegate,  epoch"
	} else {
		sqlGroup = "to_char(time_bucket, $INTERVAL), delegate"
	}
	sqlGroup = strings.Replace(sqlGroup, "$INTERVAL", "'"+timeInterval.String()+"'", -1)
	scope = scope.Group(sqlGroup)

	sqlOrder := "to_char(time_bucket, $INTERVAL)"
	sqlOrder = strings.Replace(sqlOrder, "$INTERVAL", "'"+timeInterval.String()+"'", -1)
	scope = scope.Order(sqlOrder)

	res := []model.RewardsSummary{}
	err := scope.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Import creates new rewards
func (s RewardStore) Import(records []model.BlockReward) error {
	if len(records) == 0 {
		return nil
	}

	return bulk.Import(s.db, queries.BlockRewardsImport, len(records), func(idx int) bulk.Row {
		tx := records[idx]
		return bulk.Row{
			tx.OwnerAccount,
			tx.Delegate,
			tx.Epoch,
			tx.TimeBucket,
			tx.Reward,
			tx.OwnerType,
		}
	})
}
