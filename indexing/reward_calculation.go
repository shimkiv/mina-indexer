package indexing

import (
	"errors"
	"math/big"

	log "github.com/sirupsen/logrus"

	"github.com/figment-networks/mina-indexer/model/util"
	"github.com/figment-networks/mina-indexer/store"
)

// RewardCalculation calculates rewards
func RewardCalculation(db *store.Store, data *Data) error {
	if data.CreatorFee.Float == nil || data.Block.Coinbase.Int == nil || data.Block.TransactionsFees.Int == nil || data.Block.SnarkJobsFees.Int64() == 0 {
		return nil
	}
	blockReward := data.Block.Coinbase.
		Add(data.Block.TransactionsFees).
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
		recordsMap[r.PublicKey] = *r.Weight.Float
	}

	for _, dbr := range data.DelegatorsBlockRewards {
		weight, ok := recordsMap[dbr.OwnerAccount]
		if !ok {
			err = errors.New("record is not found for " + dbr.OwnerAccount)
			log.WithError(err)
			return err
		}

		res, err := util.CalculateDelegatorReward(weight, blockReward, data.CreatorFee)
		if err != nil {
			return err
		}
		dbr.Reward = res
	}

	validatorReward, err := util.CalculateValidatorReward(blockReward, data.CreatorFee)
	if err != nil {
		return err
	}
	data.ValidatorBlockReward.Reward = validatorReward

	return nil
}
