package util

import (
	"github.com/figment-networks/mina-indexer/model/types"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestCalculateWeight(t *testing.T) {
	type args struct {
		balance            types.Amount
		totalStakedBalance types.Amount
	}
	tests := []struct {
		name    string
		args    args
		result  *big.Float
		wantErr bool
	}{
		{
			name: "successful",
			args: args{
				balance:            types.NewInt64Amount(10),
				totalStakedBalance: types.NewInt64Amount(10000),
			},
			result: big.NewFloat(0.001),
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
