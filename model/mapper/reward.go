package mapper

import (
	"errors"
	"strings"

	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/types"
)

// ValidatorBlockReward returns validator reward models references from the block data
func ValidatorBlockReward(v *model.Validator) (*model.BlockReward, error) {
	t, err := RewardTimeBucket(v.LastTime)
	if err != nil {
		return nil, err
	}
	result := model.BlockReward{
		OwnerAccount: v.PublicKey,
		OwnerType:    string(model.RewardOwnerTypeValidator),
	}
	result.TimeBucket = t
	return &result, nil
}

// DelegatorBlockRewards returns delegator reward models references from the block data
func DelegatorBlockRewards(accounts []model.LedgerEntry, block *graph.Block) ([]model.BlockReward, error) {
	result := []model.BlockReward{}
	for _, a := range accounts {
		bt := BlockTime(block)
		t, err := RewardTimeBucket(bt)
		if err != nil {
			return nil, err
		}
		// reward to be calculated next step
		dbr := model.BlockReward{
			OwnerAccount: a.PublicKey,
			Delegate:     a.Delegate,
			OwnerType:    string(model.RewardOwnerTypeDelegator),
		}
		dbr.TimeBucket = t
		result = append(result, dbr)
	}
	return result, nil
}

// ValidatorBlockReward returns validator reward models references from the block data
func FindValidatorFee(validatorEpochs []model.ValidatorEpoch, creator string) (types.Percentage, error) {
	for _, ve := range validatorEpochs {
		if ve.AccountAddress == creator {
			return ve.ValidatorFee, nil
		}
	}
	return types.Percentage{}, errors.New("validator fee not found")
}

func TransactionFees(block *graph.Block) types.Amount {
	res := types.NewInt64Amount(0)
	if block.Transactions != nil && block.Transactions.FeeTransfer != nil {
		for _, transfer := range block.Transactions.FeeTransfer {
			if strings.ToLower(transfer.Type) == model.TxTypeFeeTransfer {
				res = res.Add(types.NewAmount(transfer.Fee))
			} else if strings.ToLower(transfer.Type) == model.TxTypeCoinbaseFeeTransfer {
				res = res.Sub(types.NewAmount(transfer.Fee))
			}
		}
	}
	return res
}
