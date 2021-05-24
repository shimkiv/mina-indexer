package indexing

import (
	"errors"
	"github.com/figment-networks/mina-indexer/model/util"
	"github.com/figment-networks/mina-indexer/store"
	log "github.com/sirupsen/logrus"
	"math/big"
)

// RewardCalculation calculates rewards
func RewardCalculation(db *store.Store, data *Data) error {
	blockReward := data.Block.Coinbase.
		Mul(data.Block.TransactionsFees).
		Sub(data.Block.SnarkJobsFees)

	ledger, err := db.Staking.FindLedger(data.Block.Epoch)
	if err != nil && err != store.ErrNotFound {
		return err
	}
	if err == store.ErrNotFound {
		return nil
	}

	records, err := db.Staking.LedgerRecords(ledger.ID)
	if err != nil && err != store.ErrNotFound {
		return err
	}

	recordsMap := map[string]big.Float{}
	for _, r := range records {
		recordsMap[r.PublicKey] = r.Weight
	}

	for _, dbr := range data.DelegatorBlockRewards {
		weight, ok := recordsMap[dbr.PublicKey]
		if !ok {
			err = errors.New("record is not found for " + dbr.PublicKey)
			log.WithError(err)
			return err
		}

		res, err := util.CalculateDelegatorReward(weight, blockReward)
		if err != nil {
			return err
		}
		dbr.Reward = res
	}
	return nil
}
