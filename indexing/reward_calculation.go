package indexing

import (
	"errors"
	"math/big"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/mapper"
	"github.com/figment-networks/mina-indexer/model/types"
	"github.com/figment-networks/mina-indexer/model/util"
	"github.com/figment-networks/mina-indexer/store"
)

// RewardCalculation calculates rewards
func RewardCalculation(db *store.Store, block model.Block) error {
	if block.TransactionsFees.Int == nil || block.SnarkJobsFees.Int == nil {
		return nil
	}

	validatorEpochs, err := db.ValidatorsEpochs.GetValidatorEpochs(strconv.Itoa(block.Epoch), block.Creator)
	if err != nil && err != store.ErrNotFound {
		return err
	} else if len(validatorEpochs) == 0 {
		log.WithField("validator", block.Creator).Warn("epoch commission rate not found")
		return nil
	}

	creatorFee := validatorEpochs[0].ValidatorFee
	if err != nil {
		return err
	}
	blockReward := block.Coinbase.Add(block.TransactionsFees)
	blockReward = blockReward.Sub(block.SnarkJobsFees)

	ledger, err := db.Staking.FindLedger(block.Epoch)
	if err != nil && err != store.ErrNotFound {
		return err
	}

	delegations, err := db.Staking.FindDelegations(store.FindDelegationsParams{
		Delegate: block.Creator,
		LedgerID: &ledger.ID,
	})
	if err != nil && err != store.ErrNotFound {
		return err
	}

	firstBlockOfEpoch, err := db.Blocks.FirstBlockOfEpoch(strconv.Itoa(block.Epoch))
	if err != nil {
		return err
	}

	firstSlotOfEpoch := firstBlockOfEpoch.Slot

	if !block.Supercharged {
		err = util.CalculateWeightsNonSupercharged(delegations)
		if err != nil {
			return err
		}
	} else {
		superchargedWeighting, err := util.CalculateSuperchargedWeighting(block)
		if err != nil {
			return err
		}
		records, err := db.Staking.DelegateLedgerRecords(ledger.ID, block.Creator)
		if err != nil && err != store.ErrNotFound {
			return err
		}

		err = util.CalculateWeightsSupercharged(superchargedWeighting, delegations, records, firstSlotOfEpoch)
		if err != nil {
			return err
		}
	}

	validatorReward, err := mapper.ValidatorBlockReward(block)
	if err != nil {
		return err
	}
	reward, err := util.CalculateValidatorReward(blockReward, creatorFee)
	if err != nil {
		return err
	}
	r := new(big.Int)
	reward.Int(r)
	validatorReward.Reward = types.NewAmount(r.String())

	remainingReward := types.NewFloat(blockReward.String())
	remainingReward = remainingReward.Sub(reward)

	recordsMap := map[string]big.Float{}
	for _, r := range delegations {
		recordsMap[r.PublicKey] = *r.Weight.Float
	}

	delegatorsBlockRewards, err := mapper.DelegatorBlockRewards(delegations, block)
	if err != nil {
		return err
	}

	for i, dbr := range delegatorsBlockRewards {
		weight, ok := recordsMap[dbr.OwnerAccount]
		if !ok {
			err = errors.New("record is not found for " + dbr.OwnerAccount)
			log.WithError(err)
			return err
		}
		res, err := util.CalculateDelegatorReward(weight, remainingReward)
		if err != nil {
			return err
		}
		r = new(big.Int)
		res.Int(r)
		delegatorsBlockRewards[i].Reward = types.NewAmount(r.String())
	}

	if err := db.Rewards.Import(delegatorsBlockRewards); err != nil {
		return err
	}

	if err := db.Rewards.Import([]model.BlockReward{*validatorReward}); err != nil {
		return err
	}

	block.RewardCalculated = true
	err = db.Blocks.Update(block)
	if err != nil {
		return err
	}
	return nil
}
