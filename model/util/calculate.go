package util

import (
	"errors"
	"github.com/figment-networks/mina-indexer/model/types"
	"math/big"
)

// CalculateWeight calculates weight
func CalculateWeight(balance types.Amount, totalStakedAmount types.Amount) (big.Float, error) {
	w, ok := new(big.Float).SetString(balance.String())
	if !ok {
		return big.Float{}, errors.New("error with balance amount")
	}

	if totalStakedAmount.Int64() == 0 {
		return big.Float{}, errors.New("total staked amount can not be zero")
	}
	t, ok := new(big.Float).SetString(totalStakedAmount.String())
	if !ok {
		return big.Float{}, errors.New("error with total staked amount")
	}

	return *new(big.Float).Quo(w, t), nil
}
