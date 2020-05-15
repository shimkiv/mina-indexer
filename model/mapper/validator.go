package mapper

import (
	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/model"
)

// Validator returns a validator model constructed from the coda input
func Validator(block coda.Block) (*model.Validator, error) {
	validator := &model.Validator{
		PublicKey: block.Creator,
		Height:    blockHeight(block),
		Time:      blockTime(block),
	}
	return validator, validator.Validate()
}
