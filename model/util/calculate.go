package util

import (
	"errors"
	"math"
	"math/big"

	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/types"
)

const slotsPerEpoch = 7140

var zero = types.NewInt64Amount(0)

// CalculateWeight calculates weight
func CalculateWeight(balance types.Amount, totalStakedAmount types.Amount) (types.Float, error) {
	if totalStakedAmount.Compare(zero) == 0 {
		return types.Float{}, errors.New("total staked amount can not be zero")
	}
	w := types.NewFloat(balance.String())
	t := types.NewFloat(totalStakedAmount.String())
	return w.Quo(t), nil
}

// CalculateValidatorReward calculates validator reward
func CalculateValidatorReward(blockReward types.Amount, validatorFee types.Float) (types.Float, error) {
	vr := types.NewFloat(blockReward.String())
	fee := validatorFee.Quo(types.NewFloat("100"))
	vr = vr.Mul(fee)
	result := new(big.Int)
	vr.Int(result)
	return types.NewFloat(result.String()), nil
}

// CalculateDelegatorReward calculates delegator reward
func CalculateDelegatorReward(weight big.Float, remainingReward types.Float) (types.Float, error) {
	w := types.NewFloat(weight.String())
	return w.Mul(remainingReward), nil
}

// CalculateSuperchargedWeighting calculates supercharged weighting for given block
// supercharged weighting = 1 + (1 / (1 + transaction fees / coinbase))
func CalculateSuperchargedWeighting(block model.Block) (types.Float, error) {
	trFees := types.NewFloat(block.TransactionsFees.String())
	coinbase := types.NewFloat(block.Coinbase.String())
	denom := trFees.Quo(coinbase)
	denom = denom.Add(types.NewFloat("1"))
	enum := types.NewFloat("1")
	enum = enum.Quo(denom)
	res := enum.Add(types.NewFloat("1"))
	return res, nil
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
func CalculateWeightsSupercharged(superchargedWeighting types.Float, delegations []model.Delegation, records []model.LedgerEntry) error {
	sum := new(big.Float)
	supercharged := new(big.Float)
	balance := new(big.Float)
	stake := new(big.Float)
	weight := new(big.Float)
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

		timedWeighting, err := calculateTimedWeighting(record)
		if err != nil {
			return err
		}

		superchargedContribution, err := calculateSuperchargedContribution(superchargedWeighting, timedWeighting)
		if err != nil {
			return err
		}

		supercharged, ok = supercharged.SetString(superchargedContribution.String())
		if !ok {
			return errors.New("error with supercharged contribution")
		}
		bln, ok := balance.SetString(r.Balance.String())
		if !ok {
			return errors.New("error with balance")
		}
		stake = bln.Mul(bln, supercharged)

		delegations[i].Weight = types.NewFloat(stake.String())
		sum.Add(sum, stake)
	}

	// pool weights
	for i, r := range delegations {
		weight, ok = new(big.Float).SetString(r.Weight.String())
		if !ok {
			return errors.New("error with weight")
		}
		weight = weight.Quo(weight, sum)
		delegations[i].Weight = types.NewFloat(weight.String())
	}

	return nil
}

// calculateTimedWeighting calculates timed weighting
//
// return 1 for if entry is unlocked entire epoch,
//        0 for if entry is locked entire epoch,
//        proportion calculated based on slots if it becomes unlocked during epoch
func calculateTimedWeighting(record model.LedgerEntry) (types.Float, error) {
	if record.IsUntimed() {
		return types.NewFloat64Float(1), nil
	}

	var globalSlotStart types.Amount
	if *record.TimingVestingPeriod == 0 || record.TimingVestingIncrement.Int == nil || record.TimingVestingIncrement.Int64() == 0 {
		globalSlotStart = types.NewInt64Amount(int64(*record.TimingCliffTime))
	} else {
		globalSlotStart = record.TimingInitialMinimumBalance.Sub(record.TimingCliffAmount)
		globalSlotStart = globalSlotStart.Quo(types.NewAmount(record.TimingVestingIncrement.String()))
		globalSlotStart = globalSlotStart.Mul(types.NewInt64Amount(int64(*record.TimingVestingPeriod)))
		globalSlotStart = globalSlotStart.Add(types.NewInt64Amount(int64(*record.TimingCliffTime)))
	}
	globalSlotEnd := globalSlotStart.Add(types.NewInt64Amount(slotsPerEpoch))

	// unlockedTime: global slot at which that account's tokens will be fully unlocked
	//
	// vesting_amount = initial_min_balance - cliff_amount
	// vesting_time = vesting_amount * math.ceil(vesting_period / vesting_increment)
	// unlocked_time = cliff_time + vesting_time
	imb := types.NewAmount(record.TimingInitialMinimumBalance.String())
	ca := types.NewAmount(record.TimingCliffAmount.String())
	ct := types.NewInt64Amount(int64(*record.TimingCliffTime))
	vestingAmount := imb.Sub(ca)
	vestingPerc := math.Ceil(float64(*record.TimingVestingPeriod) / float64(record.TimingVestingIncrement.Int64()))
	vestingTime := types.NewInt64Amount(int64(vestingPerc))
	vestingTime = vestingTime.Mul(vestingAmount)
	unLockedTime := ct.Add(vestingTime)

	factor := types.NewFloat(globalSlotEnd.String()).Sub(types.NewFloat(unLockedTime.String()))
	res := types.NewFloat("1").Sub(factor.Quo(types.NewFloat64Float(slotsPerEpoch)))
	return res, nil
}

// calculateSuperchargedContribution calculates supercharged contribution
//
// supercharged contribution = ((supercharged weighting - 1) * timed weighting factor) + 1
func calculateSuperchargedContribution(superchargedWeighting, timedWeighting types.Float) (types.Float, error) {
	res := superchargedWeighting.Sub(types.NewFloat("1"))
	res = res.Mul(timedWeighting)
	res = res.Add(types.NewFloat("1"))
	return res, nil
}
