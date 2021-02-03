package mapper

import (
	"errors"
	"time"

	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/model/util"
)

var (
	errNoProtocolState   = errors.New("no protocol state")
	errNoConsensusState  = errors.New("no consensus state")
	errNoBlockchainState = errors.New("no blockchain state")
)

// BlockHeight returns a parsed block height
func BlockHeight(input *graph.Block) uint64 {
	// NOTE: graph's height starts at height=2!
	return util.MustUInt64(input.ProtocolState.ConsensusState.BlockHeight)
}

// BlockTime returns a parsed block time
func BlockTime(input *graph.Block) time.Time {
	return util.MustTime(input.ProtocolState.BlockchainState.Date)
}

func blockCheck(input *graph.Block) error {
	if input.ProtocolState == nil {
		return errNoProtocolState
	}
	if input.ProtocolState.ConsensusState == nil {
		return errNoConsensusState
	}
	if input.ProtocolState.BlockchainState == nil {
		return errNoBlockchainState
	}
	return nil
}
