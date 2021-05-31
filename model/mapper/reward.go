package mapper

import (
	"errors"

	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/types"
)

// ValidatorBlockReward returns validator reward models references from the block data
func ValidatorBlockReward(v *model.Validator) (*model.BlockReward, error) {
	result := model.BlockReward{
		OwnerAccount: v.PublicKey,
		BlockHeight:  v.LastHeight,
		BlockTime:    v.LastTime,
		OwnerType:    string(model.RewardOwnerTypeValidator),
	}
	return &result, nil
}

// DelegatorBlockRewards returns delegator reward models references from the block data
func DelegatorBlockRewards(accounts []model.Account) ([]model.BlockReward, error) {
	result := []model.BlockReward{}
	for _, a := range accounts {
		// reward to be calculated next step
		dbr := model.BlockReward{
			OwnerAccount: a.PublicKey,
			BlockHeight:  a.LastHeight,
			BlockTime:    a.LastTime,
			OwnerType:    string(model.RewardOwnerTypeDelegator),
		}
		if a.Delegate != nil {
			dbr.Delegate = *a.Delegate
		} else {
			// self delegation
			dbr.Delegate = a.PublicKey
		}
		result = append(result, dbr)
	}

	return result, nil
}

// ValidatorBlockReward returns validator reward models references from the block data
func FindValidatorFee(validatorEpochs []model.ValidatorEpoch, creator string) (types.Percentage, error) {
	for _, ve := range validatorEpochs {
		if ve.AccountId == creator {
			return ve.ValidatorFee, nil
		}
	}
	return types.Percentage{}, errors.New("validator fee not found")
}

func CoinbaseReward(block *graph.Block) types.Amount {
	if block.Transactions != nil {
		return types.NewAmount(block.Transactions.Coinbase)
	}
	return types.NewAmount("")
}

func TransactionFees(block *graph.Block) types.Amount {
	res := types.NewInt64Amount(0)
	if block.Transactions != nil && block.Transactions.FeeTransfer != nil {
		for _, transfer := range block.Transactions.FeeTransfer {
			res = res.Add(types.NewAmount(transfer.Fee))
		}
	}
	return res
}
