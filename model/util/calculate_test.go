package util

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/types"
)

func TestCalculateWeight(t *testing.T) {
	type args struct {
		balance            types.Amount
		totalStakedBalance types.Amount
	}
	tests := []struct {
		name    string
		args    args
		result  types.Float
		wantErr bool
	}{
		{
			name: "successful",
			args: args{
				balance:            types.NewInt64Amount(10),
				totalStakedBalance: types.NewInt64Amount(10000),
			},
			result: types.NewFloat("0.001"),
		},
		{
			name: "error case stake value",
			args: args{
				balance:            types.NewInt64Amount(10),
				totalStakedBalance: types.NewInt64Amount(0),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := CalculateWeight(tt.args.balance, tt.args.totalStakedBalance)
			if err != nil {
				assert.True(t, tt.wantErr)
			} else {
				assert.Equal(t, res.String(), tt.result.String())
			}
		})
	}
}

func TestCalculateDelegatorReward(t *testing.T) {
	w, _ := new(big.Float).SetString("0.3")
	type args struct {
		weight          big.Float
		remainingReward types.Float
	}
	tests := []struct {
		name    string
		args    args
		result  types.Float
		wantErr bool
	}{
		{
			name: "successful",
			args: args{
				weight:          *w,
				remainingReward: types.NewFloat("95"),
			},
			result: types.NewFloat("28.5"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := CalculateDelegatorReward(tt.args.weight, tt.args.remainingReward)
			if err != nil {
				assert.True(t, tt.wantErr)
			} else {
				assert.Equal(t, res.String(), tt.result.String())
			}
		})
	}
}

func TestCalculateValidatorReward(t *testing.T) {
	type args struct {
		blockReward  types.Amount
		validatorFee types.Float
	}
	tests := []struct {
		name    string
		args    args
		result  types.Amount
		wantErr bool
	}{
		{
			name: "successful",
			args: args{
				blockReward:  types.NewInt64Amount(100),
				validatorFee: types.NewFloat("5"),
			},
			result: types.NewAmount("5"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := CalculateValidatorReward(tt.args.blockReward, tt.args.validatorFee)
			if err != nil {
				assert.True(t, tt.wantErr)
			} else {
				assert.Equal(t, res.String(), tt.result.String())
			}
		})
	}
}

func TestCalculateSuperchargedWeighting(t *testing.T) {
	type args struct {
		block model.Block
	}
	tests := []struct {
		name    string
		args    args
		result  types.Float
		wantErr bool
	}{
		{
			name: "successful zero transaction fee",
			args: args{
				block: model.Block{
					Coinbase:         types.NewAmount("200"),
					TransactionsFees: types.NewAmount("0"),
				},
			},
			result: types.NewFloat("2"),
		},
		{
			name: "successful same amount transaction fee and coinbase",
			args: args{
				block: model.Block{
					Coinbase:         types.NewAmount("200"),
					TransactionsFees: types.NewAmount("200"),
				},
			},
			result: types.NewFloat("1.5"),
		},
		{
			name: "successful low transaction fee",
			args: args{
				block: model.Block{
					Coinbase:         types.NewAmount("200"),
					TransactionsFees: types.NewAmount("5"),
				},
			},
			result: types.NewFloat("1.975609756"),
		},
		{
			name: "successful high transaction fee",
			args: args{
				block: model.Block{
					Coinbase:         types.NewAmount("200"),
					TransactionsFees: types.NewAmount("10000"),
				},
			},
			result: types.NewFloat("1.019607843"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := CalculateSuperchargedWeighting(tt.args.block)
			if err != nil {
				assert.True(t, tt.wantErr)
			} else {
				assert.Equal(t, res.String(), tt.result.String())
			}
		})
	}
}

func TestCalculateWeightsSupercharged(t *testing.T) {
	ct := 86400
	vp := 1

	records := []model.LedgerEntry{
		{
			PublicKey: "one",
			Balance:   types.NewAmount("20000"),
		},
		{
			PublicKey: "two",
			Balance:   types.NewAmount("50000"),
		},
		{
			PublicKey:                   "three",
			Balance:                     types.NewAmount("30000"),
			TimingInitialMinimumBalance: types.NewAmount("230400"),
			TimingCliffAmount:           types.NewAmount("230400"),
			TimingCliffTime:             &ct,
			TimingVestingPeriod:         &vp,
			TimingVestingIncrement:      types.NewAmount("100"),
		},
	}

	delegations := []model.Delegation{
		{
			PublicKey: "one",
			Balance:   types.NewAmount("20000"),
		},
		{
			PublicKey: "two",
			Balance:   types.NewAmount("50000"),
		},
		{
			PublicKey: "three",
			Balance:   types.NewAmount("30000"),
		},
	}

	type args struct {
		superchargedContribution types.Float
		delegations              []model.Delegation
		records                  []model.LedgerEntry
		firstSlotOfEpoch         int
	}
	tests := []struct {
		name    string
		args    args
		result  []string
		wantErr bool
	}{
		{
			name: "successful with timing but locked entire epoch",
			args: args{
				superchargedContribution: types.NewFloat("1.98765"),
				delegations:              delegations,
				records:                  records,
			},
			result: []string{"0.2350364057", "0.5875910143", "0.17737258"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delegations := tt.args.delegations
			err := CalculateWeightsSupercharged(tt.args.superchargedContribution, delegations, tt.args.records)
			if err != nil {
				assert.True(t, tt.wantErr)
			} else {
				for i, _ := range records {
					assert.Equal(t, tt.result[i], delegations[i].Weight.String())
				}
			}
		})
	}
}
