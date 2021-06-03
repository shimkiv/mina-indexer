package util

import (
	"errors"
	"math"
	"math/big"

	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/types"
)

const slotsPerEpoch = 7140

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
func CalculateValidatorReward(blockReward types.Amount, validatorFee types.Percentage) (types.Percentage, error) {
	vr, ok := new(big.Float).SetString(blockReward.String())
	if !ok {
		return types.Percentage{}, errors.New("error with block reward amount")
	}

	fee, ok := new(big.Float).SetString(validatorFee.String())
	if !ok {
		return types.Percentage{}, errors.New("error with validator fee")
	}
	fee.Quo(fee, big.NewFloat(100))
	vr.Mul(vr, fee)
	if !ok {
		return types.Percentage{}, errors.New("error with validator reward amount")
	}

	result := new(big.Int)
	vr.Int(result)
	return types.NewPercentage(result.String()), nil
}

// CalculateDelegatorReward calculates delegator reward
func CalculateDelegatorReward(weight big.Float, blockReward types.Amount, validatorFee types.Percentage) (types.Percentage, error) {
	br, ok := new(big.Float).SetString(blockReward.String())
	if !ok {
		return types.Percentage{}, errors.New("error with stake amount")
	}
	remaining := big.NewFloat(0)
	remaining.Sub(big.NewFloat(100), validatorFee.Float)
	remaining.Quo(remaining, big.NewFloat(100))
	br = br.Mul(br, remaining)
	if !ok {
		return types.Percentage{}, errors.New("error with stake amount")
	}
	res := new(big.Float)
	res.Mul(br, &weight)
	return types.NewPercentage(res.String()), nil
}

// CalculateSuperchargedWeighting calculates supercharged weighting
func CalculateSuperchargedWeighting(block model.Block) (types.Percentage, error) {
	trFees, ok := new(big.Float).SetString(block.TransactionsFees.String())
	if !ok {
		return types.Percentage{}, errors.New("error with transaction fees")
	}
	coinbase, ok := new(big.Float).SetString(block.Coinbase.String())
	if !ok {
		return types.Percentage{}, errors.New("error with coinbase")
	}
	denom := trFees.Quo(trFees, coinbase)
	denom = denom.Add(denom, new(big.Float).SetFloat64(1))
	enum := new(big.Float).SetFloat64(1)
	enum = enum.Quo(enum, denom)
	res := enum.Add(enum, new(big.Float).SetFloat64(1))
	return types.NewPercentage(res.String()), nil
}

// CalculateWeightsNonSupercharged calculates weights when block reward is not doubled for supercharged
func CalculateWeightsNonSupercharged(StakedAmount types.Amount, records []model.LedgerEntry) error {
	for _, r := range records {
		w, err := CalculateWeight(r.Balance, StakedAmount)
		if err != nil {
			return err
		}
		r.Weight = w
	}
	return nil
}

// CalculateWeightsSupercharged calculates weights when block reward is doubled for supercharged
func CalculateWeightsSupercharged(superchargedWeighting types.Percentage, records []model.LedgerEntry, firstSlotOfEpoch int) error {
	sum := new(big.Float)
	sc := new(big.Float)
	bln := new(big.Float)
	stk := new(big.Float)
	w := new(big.Float)
	var ok bool

	// pool stakes
	for i, r := range records {
		timedWeighting, err := calculateTimedWeighting(r, firstSlotOfEpoch)
		if err != nil {
			return err
		}

		superchargedContribution, err := calculateSuperchargedContribution(superchargedWeighting, timedWeighting)
		if err != nil {
			return err
		}

		if r.LockedTokens.Int64() > 0 {
			stk, ok = bln.SetString(r.Balance.String())
			if !ok {
				return errors.New("error with balance")
			}
		} else {
			sc, ok = sc.SetString(superchargedContribution.String())
			if !ok {
				return errors.New("error with supercharged contribution")
			}
			bln, ok := bln.SetString(r.Balance.String())
			if !ok {
				return errors.New("error with balance")
			}
			stk = bln.Mul(bln, sc)
		}
		records[i].Weight = types.NewPercentage(stk.String())
		sum.Add(sum, stk)
	}

	// pool weights
	for i, r := range records {
		w, ok = new(big.Float).SetString(r.Weight.String())
		if !ok {
			return errors.New("error with weight")
		}
		w = w.Quo(w, sum)
		records[i].Weight = types.NewPercentage(w.String())
	}

	return nil
}

// calculateTimedWeighting calculates timed weighting
func calculateTimedWeighting(record model.LedgerEntry, firstSlotOfEpoch int) (types.Percentage, error) {
	if record.IsUntimed() && record.LockedTokens.Int64() == 0 {
		return types.NewFloat64Percentage(1), nil
	}

	if record.IsUntimed() && record.LockedTokens.Int64() > 0 {
		return types.NewFloat64Percentage(0), nil
	}

	imb, ok := new(big.Int).SetString(record.TimingInitialMinimumBalance.String(), 10)
	if !ok {
		return types.Percentage{}, errors.New("error with initial minimum balance")
	}
	ca, ok := new(big.Int).SetString(record.TimingCliffAmount.String(), 10)
	if !ok {
		return types.Percentage{}, errors.New("error with cliff amount")
	}
	ct := new(big.Int).SetInt64(int64(*record.TimingCliffTime))

	vestingAmount := new(big.Int)
	vestingAmount.Sub(imb, ca)
	vestingPerc := math.Ceil(float64(*record.TimingVestingPeriod) / float64(*record.TimingVestingIncrement))
	vestingTime := new(big.Int).SetInt64(int64(vestingPerc))
	vestingTime.Mul(vestingTime, vestingAmount)
	unLockedTime := new(big.Int)
	unLockedTime.Add(ct, vestingTime)

	factor := float64(slotsPerEpoch-(unLockedTime.Int64()-int64(firstSlotOfEpoch))) / float64(slotsPerEpoch)
	return types.NewFloat64Percentage(factor), nil
}

// calculateSuperchargedContribution calculates supercharged contribution
func calculateSuperchargedContribution(superchargedWeighting, timedWeighting types.Percentage) (types.Percentage, error) {
	sw, ok := new(big.Float).SetString(superchargedWeighting.String())
	if !ok {
		return types.Percentage{}, errors.New("error with supercharged weighting")
	}
	tw, ok := new(big.Float).SetString(timedWeighting.String())
	if !ok {
		return types.Percentage{}, errors.New("error with timed weighting")
	}
	res := sw.Sub(sw, big.NewFloat(1))
	res = res.Mul(sw, tw)
	res = res.Add(sw, big.NewFloat(1))
	return types.NewPercentage(res.String()), nil
}
