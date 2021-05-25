package store

import (
	"strings"
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/store/queries"
)

// DelegatorRewardStore handles operations on delegator rewards
type DelegatorRewardStore struct {
	baseStore
}

// FetchRewardsByInterval fetches rewards by interval
func (s *DelegatorRewardStore) FetchRewardsByInterval(publicKey string, delegate string, from time.Time, to time.Time, timeInterval model.TimeInterval) (model.RewardsSummary, error) {
	var res model.RewardsSummary
	q := strings.Replace(queries.DelegatorsRewards, "$INTERVAL", "'"+timeInterval.String()+"'", -1)
	var err error
	if delegate == "" {
		q = strings.Replace(q, "AND delegate = ?", "", -1)
		err = s.db.Raw(q, publicKey, from, to).Scan(&res).Error
	} else {
		err = s.db.Raw(q, publicKey, delegate, from, to).Scan(&res).Error
	}
	if err != nil {
		return res, err
	}
	return res, nil
}

// Import creates new delegator rewards
func (s DelegatorRewardStore) Import(records []model.DelegatorBlockReward) error {
	if len(records) == 0 {
		return nil
	}

	return bulk.Import(s.db, queries.DelegatorBlockRewardsImport, len(records), func(idx int) bulk.Row {
		tx := records[idx]
		return bulk.Row{
			tx.PublicKey,
			tx.Delegate,
			tx.BlockHeight,
			tx.BlockTime,
			tx.Reward,
		}
	})
}
