package mapper

import (
	"errors"
	"time"

	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/model/util"
)

var (
	errNoProtocolState   = errors.New("no protocol state")
	errNoConsensusState  = errors.New("no consensus state")
	errNoBlockchainState = errors.New("no blockchain state")
)

func blockCheck(input coda.Block) error {
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

func blockHeight(input coda.Block) int64 {
	return util.MustInt64(input.ProtocolState.ConsensusState.BlockHeight)
}

func blockTime(input coda.Block) time.Time {
	return util.MustTime(input.ProtocolState.BlockchainState.Date)
}
