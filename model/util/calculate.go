package util

import (
	"errors"
	"math/big"

	"github.com/figment-networks/mina-indexer/model/types"
)

// CalculateWeight calculates weight
func CalculateWeight(balance types.Amount, totalStakedAmount types.Amount) (types.Percentage, error) {
	w, ok := new(big.Float).SetString(balance.String())
	if !ok {
		return types.Percentage{}, errors.New("error with balance amount")
	}

	if totalStakedAmount.Int64() == 0 {
		return types.Percentage{}, errors.New("total staked amount can not be zero")
	}
	t, ok := new(big.Float).SetString(totalStakedAmount.String())
	if !ok {
		return types.Percentage{}, errors.New("error with total staked amount")
	}

	return types.NewPercentage(new(big.Float).Quo(w, t).String()), nil
}

// CalculateValidatorReward calculates validator reward
func CalculateValidatorReward(blockReward types.Amount) (types.Amount, error) {
	vr, ok := new(big.Float).SetString(blockReward.String())
	if !ok {
		return types.Amount{}, errors.New("error with block reward amount")
	}

	// %5 validator fee
	// TODO: fetch validators fee from https://api.staketab.com/mina/get_providers
	vr = vr.Mul(vr, big.NewFloat(0.05))
	if !ok {
		return types.Amount{}, errors.New("error with validator reward amount")
	}

	return types.NewAmount(vr.String()), nil
}

// CalculateDelegatorReward calculates delegator reward
func CalculateDelegatorReward(weight big.Float, blockReward types.Amount) (types.Amount, error) {
	br, ok := new(big.Float).SetString(blockReward.String())
	if !ok {
		return types.Amount{}, errors.New("error with stake amount")
	}

	// %5 validator fee
	br = br.Mul(br, big.NewFloat(0.95))
	if !ok {
		return types.Amount{}, errors.New("error with stake amount")
	}

	res := new(big.Float)
	res.Mul(br, &weight)
	return types.NewAmount(res.String()), nil
}
