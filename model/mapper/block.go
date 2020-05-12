package mapper

import (
	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/model"
	"github.com/figment-networks/coda-indexer/model/util"
)

// Block returns a block model constructed from the coda input
func Block(input coda.Block) (*model.Block, error) {
	blockTime, err := util.ParseTime(input.ProtocolState.BlockchainState.Date)
	if err != nil {
		return nil, err
	}

	blockHeight, err := util.ParseInt64(input.ProtocolState.ConsensusState.BlockHeight)
	if err != nil {
		return nil, err
	}

	coinbase, err := util.ParseInt64(input.Transactions.Coinbase)
	if err != nil {
		return nil, err
	}

	block := &model.Block{
		Height:            blockHeight,
		Time:              *blockTime,
		Creator:           input.Creator,
		Hash:              input.StateHash,
		ParentHash:        input.ProtocolState.PreviousStateHash,
		LedgerHash:        input.ProtocolState.BlockchainState.SnarkedLedgerHash,
		TransactionsCount: len(input.Transactions.UserCommands),
		Coinbase:          coinbase,
	}

	return block, nil
}
