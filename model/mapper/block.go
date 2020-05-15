package mapper

import (
	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/model"
	"github.com/figment-networks/coda-indexer/model/util"
)

// Block returns a block model constructed from the coda input
func Block(input coda.Block) (*model.Block, error) {
	if err := blockCheck(input); err != nil {
		return nil, err
	}

	block := &model.Block{
		Height:            blockHeight(input),
		Time:              blockTime(input),
		Creator:           input.Creator,
		Hash:              input.StateHash,
		ParentHash:        input.ProtocolState.PreviousStateHash,
		LedgerHash:        input.ProtocolState.BlockchainState.SnarkedLedgerHash,
		TransactionsCount: len(input.Transactions.UserCommands),
		Coinbase:          util.MustInt64(input.Transactions.Coinbase),
	}

	return block, block.Validate()
}
