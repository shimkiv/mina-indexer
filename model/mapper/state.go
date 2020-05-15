package mapper

import (
	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/model"
	"github.com/figment-networks/coda-indexer/model/util"
)

// State returns a state model constructed from the coda input
func State(input coda.Block) (*model.State, error) {
	state := &model.State{
		Height:        blockHeight(input),
		TotalCurrency: util.MustInt64(input.ProtocolState.ConsensusState.TotalCurrency),
		LastVFROutput: input.ProtocolState.ConsensusState.LastVrfOutput,
		Epoch:         util.MustInt64(input.ProtocolState.ConsensusState.Epoch),
		EpochCount:    util.MustInt64(input.ProtocolState.ConsensusState.EpochCount),
	}
	return state, state.Validate()
}
