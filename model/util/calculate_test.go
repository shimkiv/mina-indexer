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
		result  types.Percentage
		wantErr bool
	}{
		{
			name: "successful",
			args: args{
				balance:            types.NewInt64Amount(10),
				totalStakedBalance: types.NewInt64Amount(10000),
			},
			result: types.NewPercentage("0.001"),
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
		weight       big.Float
		blockReward  types.Amount
		validatorFee types.Percentage
	}
	tests := []struct {
		name    string
		args    args
		result  types.Percentage
		wantErr bool
	}{
		{
			name: "successful",
			args: args{
				weight:       *w,
				blockReward:  types.NewInt64Amount(100),
				validatorFee: types.NewPercentage("5"),
			},
			result: types.NewPercentage("28.5"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := CalculateDelegatorReward(tt.args.weight, tt.args.blockReward, tt.args.validatorFee)
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
		validatorFee types.Percentage
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
				validatorFee: types.NewPercentage("5"),
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
		result  types.Percentage
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
			result: types.NewPercentage("2"),
		},
		{
			name: "successful same amount transaction fee and coinbase",
			args: args{
				block: model.Block{
					Coinbase:         types.NewAmount("200"),
					TransactionsFees: types.NewAmount("200"),
				},
			},
			result: types.NewPercentage("1.5"),
		},
		{
			name: "successful low transaction fee",
			args: args{
				block: model.Block{
					Coinbase:         types.NewAmount("200"),
					TransactionsFees: types.NewAmount("5"),
				},
			},
			result: types.NewPercentage("1.975609756"),
		},
		{
			name: "successful high transaction fee",
			args: args{
				block: model.Block{
					Coinbase:         types.NewAmount("200"),
					TransactionsFees: types.NewAmount("10000"),
				},
			},
			result: types.NewPercentage("1.019607843"),
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

// TODO: unit test for setting heights