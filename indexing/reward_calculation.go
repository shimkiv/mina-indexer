package indexing

import (
	"errors"
	"math/big"

	log "github.com/sirupsen/logrus"

	"github.com/figment-networks/mina-indexer/model/mapper"
	"github.com/figment-networks/mina-indexer/model/util"
	"github.com/figment-networks/mina-indexer/store"
)

// RewardCalculation calculates rewards
func RewardCalculation(db *store.Store, data *Data) error {
	if data.CreatorAccount == nil || data.Block.Coinbase.Int == nil || data.Block.TransactionsFees.Int == nil || data.Block.SnarkJobsFees.Int64() == 0 {
		return nil
	}

	creatorFee, err := mapper.FindValidatorFee(data.ValidatorEpochs, data.CreatorAccount.PublicKey)
	if err != nil {
		return err
	}
	blockReward := data.Block.Coinbase.Add(data.Block.TransactionsFees)
	blockReward = blockReward.Sub(data.Block.SnarkJobsFees)

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

	if !data.CreatorAccount.Supercharged {
		err = util.CalculateWeightsNonSupercharged(ledger.StakedAmount, records)
		if err != nil {
			return err
		}
	} else {
		superchargedWeighting, err := util.CalculateSuperchargedWeighting(*data.Block)
		if err != nil {
			return err
		}
		err = util.CalculateWeightsSupercharged(superchargedWeighting, records, data.FirstSlotOfEpoch)
		if err != nil {
			return err
		}
	}

	// update db for weights
	err = db.Staking.CreateLedgerEntries(records)
	if err != nil {
		return err
	}

	recordsMap := map[string]big.Float{}
	for _, r := range records {
		recordsMap[r.PublicKey] = *r.Weight.Float
	}

	for i, dbr := range data.DelegatorsBlockRewards {
		weight, ok := recordsMap[dbr.OwnerAccount]
		if !ok {
			err = errors.New("record is not found for " + dbr.OwnerAccount)
			log.WithError(err)
			return err
		}
		res, err := util.CalculateDelegatorReward(weight, blockReward, creatorFee)
		if err != nil {
			return err
		}
		data.DelegatorsBlockRewards[i].Reward = res
	}

	validatorReward, err := util.CalculateValidatorReward(blockReward, creatorFee)
	if err != nil {
		return err
	}
	data.ValidatorBlockReward.Reward = validatorReward

	return nil
}
