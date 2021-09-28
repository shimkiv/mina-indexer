package model

import (
	"testing"

	"github.com/figment-networks/mina-indexer/model/types"
	"github.com/stretchr/testify/assert"
)

func TestBlockEligibleForRewardsCalculation(t *testing.T) {
	examples := []struct {
		block    Block
		expected bool
	}{
		{
			expected: false,
			block:    Block{},
		},
		{
			expected: false,
			block: Block{
				Coinbase: types.NewAmount("100"),
			},
		},
		{
			expected: false,
			block: Block{
				Coinbase:         types.NewAmount("100"),
				TransactionsFees: types.NewAmount("1"),
			},
		},
		{
			expected: true,
			block: Block{
				Coinbase:         types.NewAmount("100"),
				TransactionsFees: types.NewAmount("1"),
				SnarkJobsFees:    types.NewAmount("2"),
			},
		},
	}

	for _, example := range examples {
		assert.Equal(t, example.expected, example.block.EligibleForRewardCalculation())
	}
}
