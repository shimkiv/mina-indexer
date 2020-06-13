package mapper

import (
	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/model"
)

// Validator returns a validator model constructed from the coda input
func Validator(block *coda.Block) (*model.Validator, error) {
	height := BlockHeight(block)
	time := BlockTime(block)

	validator := &model.Validator{
		Account:        block.Creator,
		StartHeight:    height,
		StartTime:      time,
		LastHeight:     height,
		LastTime:       time,
		BlocksProposed: 0, // TODO
		BlocksCreated:  1, // Gets recalculated on update
	}
	return validator, validator.Validate()
}
