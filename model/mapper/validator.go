package mapper

import (
	"time"

	"github.com/figment-networks/mina-indexer/client/archive"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/types"
)

// Validator returns a validator model constructed from the coda input
func Validator(block *archive.Block) (*model.Validator, error) {
	ts := time.Unix(block.Timestamp/1000, 0)

	validator := &model.Validator{
		PublicKey:        block.Creator,
		StartHeight:      block.Height,
		StartTime:        ts,
		LastHeight:       block.Height,
		LastTime:         ts,
		Stake:            types.NewInt64Amount(0),
		DelegatedBalance: types.NewInt64Amount(0),
		BlocksCreated:    1,
		BlocksProposed:   1,
	}

	return validator, validator.Validate()
}
