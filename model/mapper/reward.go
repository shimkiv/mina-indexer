package mapper

import (
	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/model/types"
)

func CoinbaseReward(block *graph.Block) types.Amount {
	if block.Transactions != nil {
		return types.NewAmount(block.Transactions.Coinbase)
	}
	return types.NewAmount("")
}

func TransactionFees(block *graph.Block) types.Amount {
	res := types.NewInt64Amount(0)
	if block.Transactions != nil && block.Transactions.FeeTransfer != nil {
		for _, transfer := range block.Transactions.FeeTransfer {
			res = res.Add(types.NewAmount(transfer.Fee))
		}
	}
	return res
}
