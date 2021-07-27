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
func (s *RewardStore) FetchRewardsByInterval(ownerAccount string, delegate string, from time.Time, to time.Time, timeInterval model.TimeInterval, rewardOwnerType model.RewardOwnerType) ([]model.RewardsSummary, error) {
	slt := "to_char(time_bucket, $INTERVAL) AS interval,  delegate,  epoch,  SUM(reward) AS amount"
	slt = strings.Replace(slt, "$INTERVAL", "'"+timeInterval.String()+"'", -1)

	scope := s.db.Select(slt).Table("block_rewards")

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
		scope = scope.Where("time_bucket > ?", from)
	}
	if !to.IsZero() {
		scope = scope.Where("time_bucket < ?", to)
	}

	grp := "to_char(time_bucket, $INTERVAL), delegate,  epoch"
	grp = strings.Replace(grp, "$INTERVAL", "'"+timeInterval.String()+"'", -1)
	scope = scope.Group(grp)

	ord := "to_char(time_bucket, $INTERVAL)"
	ord = strings.Replace(ord, "$INTERVAL", "'"+timeInterval.String()+"'", -1)
	scope = scope.Order(ord)

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
