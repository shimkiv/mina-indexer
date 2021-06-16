package indexing

import (
	"errors"
	"strconv"

	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/util"
	"github.com/figment-networks/mina-indexer/store"
)

// RewardCalculation calculates rewards
func RewardCalculation(db *store.Store, block model.Block) error {
	if block.Coinbase.Int == nil || block.TransactionsFees.Int == nil || block.SnarkJobsFees.Int == nil {
		return nil
	}

	validatorEpochs, err := db.ValidatorsEpochs.GetValidatorEpochs(strconv.Itoa(block.Epoch), block.Creator)
	if err != nil && err != store.ErrNotFound {
		return err
	} else if len(validatorEpochs) == 0 {
		return errors.New("validator fee for epoch not found")
	}

	//creatorFee := validatorEpochs[0].ValidatorFee
	_ = validatorEpochs[0].ValidatorFee
	if err != nil {
		return err
	}
	blockReward := block.Coinbase.Add(block.TransactionsFees)
	blockReward = blockReward.Sub(block.SnarkJobsFees)

	ledger, err := db.Staking.FindLedger(block.Epoch)
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

	firstBlockOfEpoch, err := db.Blocks.FirstBlockOfEpoch(strconv.Itoa(block.Epoch))
	if err != nil {
		if err != store.ErrNotFound {
			return err
		}
	} else if firstBlockOfEpoch == nil {
		return errors.New("first block of epoch is not found")
	}

	firstSlotOfEpoch := firstBlockOfEpoch.Slot

	if !block.Supercharged {
		err = util.CalculateWeightsNonSupercharged(ledger.StakedAmount, records)
		if err != nil {
			return err
		}
	} else {
		superchargedWeighting, err := util.CalculateSuperchargedWeighting(block)
		if err != nil {
			return err
		}
		err = util.CalculateWeightsSupercharged(superchargedWeighting, records, firstSlotOfEpoch)
		if err != nil {
			return err
		}
	}

	// update db for weights

	// todo remove
	err = db.Staking.CreateLedgerEntries(records)
	if err != nil {
		return err
	}

	// todo: remove weight from ledger_entries
	/*
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
	*/
	return nil
}
