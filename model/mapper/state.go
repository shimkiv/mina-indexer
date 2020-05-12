package mapper

import (
	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/model"
	"github.com/figment-networks/coda-indexer/model/util"
)

// State returns a state model constructed from the coda input
func State(block coda.Block) (*model.State, error) {
	consensus := block.ProtocolState.ConsensusState

	height, err := util.ParseInt64(consensus.BlockHeight)
	if err != nil {
		return nil, err
	}

	currency, err := util.ParseInt64(consensus.TotalCurrency)
	if err != nil {
		return nil, err
	}

	state := &model.State{
		Height:        height,
		TotalCurrency: currency,
	}

	if err := state.Validate(); err != nil {
		return nil, err
	}

	return state, nil
}
