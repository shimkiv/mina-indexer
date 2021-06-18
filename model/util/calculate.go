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

// CalculateSuperchargedWeighting calculates supercharged weighting for given block
// supercharged weighting = 1 + (1 / (1 + transaction fees / coinbase))
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
	denom.Add(denom, new(big.Float).SetFloat64(1))
	enum := new(big.Float).SetFloat64(1)
	enum.Quo(enum, denom)
	res := enum.Add(enum, new(big.Float).SetFloat64(1))
	return types.NewPercentage(res.String()), nil
}

// CalculateWeightsNonSupercharged calculates weights when block reward is not doubled for supercharged
func CalculateWeightsNonSupercharged(delegations []model.Delegation) error {
	stakedAmount := new(big.Int)
	bln := new(big.Int)
	var ok bool
	for _, r := range delegations {
		bln, ok = bln.SetString(r.Balance.String(), 10)
		if !ok {
			return errors.New("error with balance")
		}
		stakedAmount.Add(stakedAmount, bln)
	}
	for i, r := range delegations {
		w, err := CalculateWeight(r.Balance, types.NewAmount(stakedAmount.String()))
		if err != nil {
			return err
		}
		delegations[i].Weight = w
	}
	return nil
}

// CalculateWeightsSupercharged calculates weights when block reward is doubled for supercharged
func CalculateWeightsSupercharged(superchargedWeighting types.Percentage, delegations []model.Delegation, records []model.LedgerEntry, firstSlotOfEpoch int) error {
	sum := new(big.Float)
	sc := new(big.Float)
	bln := new(big.Float)
	stk := new(big.Float)
	w := new(big.Float)
	var ok bool

	recordsMap := map[string]model.LedgerEntry{}
	for _, r := range records {
		recordsMap[r.PublicKey] = r
	}

	// pool stakes
	for i, r := range delegations {
		record, ok := recordsMap[r.PublicKey]
		if !ok {
			return errors.New("ledger record not found")
		}

		timedWeighting, err := calculateTimedWeighting(record, firstSlotOfEpoch)
		if err != nil {
			return err
		}

		superchargedContribution, err := calculateSuperchargedContribution(superchargedWeighting, timedWeighting)
		if err != nil {
			return err
		}

		sc, ok = sc.SetString(superchargedContribution.String())
		if !ok {
			return errors.New("error with supercharged contribution")
		}
		bln, ok := bln.SetString(r.Balance.String())
		if !ok {
			return errors.New("error with balance")
		}
		stk = bln.Mul(bln, sc)

		delegations[i].Weight = types.NewPercentage(stk.String())
		sum.Add(sum, stk)
	}

	// pool weights
	for i, r := range delegations {
		w, ok = new(big.Float).SetString(r.Weight.String())
		if !ok {
			return errors.New("error with weight")
		}
		w = w.Quo(w, sum)
		delegations[i].Weight = types.NewPercentage(w.String())
	}

	return nil
}

// calculateTimedWeighting calculates timed weighting
//
// return 0 for if entry is unlocked entire epoch,
//        1 for if entry is locked entire epoch,
//        proportion calculated based on slots if it becomes unlocked during epoch
func calculateTimedWeighting(record model.LedgerEntry, firstSlotOfEpoch int) (types.Percentage, error) {
	if record.IsUntimed() {
		return types.NewFloat64Percentage(1), nil
	}

	// unlockedTime: global slot at which that account's tokens will be fully unlocked
	//
	// vesting_amount = initial_min_balance - cliff_amount
	// vesting_time = vesting_amount * math.ceil(vesting_period / vesting_increment)
	// unlocked_time = cliff_time + vesting_time
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
	if factor < 0 {
		// means it will be unlocked after this epoch, so it is locked in this one
		return types.NewFloat64Percentage(0), nil
	}
	return types.NewFloat64Percentage(factor), nil
}

// calculateSuperchargedContribution calculates supercharged contribution
//
// supercharged contribution = ((supercharged weighting - 1) * timed weighting factor) + 1
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
	res.Mul(sw, tw)
	res.Add(sw, big.NewFloat(1))
	return types.NewPercentage(res.String()), nil
}
