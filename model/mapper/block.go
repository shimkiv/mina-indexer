package mapper

import (
	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/model"
	"github.com/figment-networks/coda-indexer/model/util"
)

// Block returns a block model constructed from the coda input
func Block(input *coda.Block) (*model.Block, error) {
	if err := blockCheck(input); err != nil {
		return nil, err
	}

	block := &model.Block{
		Height:            BlockHeight(input),
		Time:              BlockTime(input),
		Creator:           input.Creator,
		Hash:              input.StateHash,
		ParentHash:        input.ProtocolState.PreviousStateHash,
		LedgerHash:        input.ProtocolState.BlockchainState.SnarkedLedgerHash,
		TransactionsCount: len(input.Transactions.UserCommands),
		FeeTransfersCount: len(input.Transactions.FeeTransfer),
		Coinbase:          util.MustUInt64(input.Transactions.Coinbase),
		TotalCurrency:     util.MustUInt64(input.ProtocolState.ConsensusState.TotalCurrency),
		Epoch:             util.MustInt64(input.ProtocolState.ConsensusState.Epoch),
		Slot:              util.MustInt64(input.ProtocolState.ConsensusState.Slot),
		SnarkJobsCount:    len(input.SnarkJobs),
	}

	snarkers := map[string]bool{}
	for _, j := range input.SnarkJobs {
		if !snarkers[j.Prover] {
			snarkers[j.Prover] = true
			block.SnarkersCount++
		}
	}

	return block, block.Validate()
}
