package indexing

import (
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/types"
)

// Data contains all the records processed for a height
type Data struct {
	Block                 *model.Block
	Validator             *model.Validator
	ValidatorBlockReward  *model.BlockReward
	CreatorFee            types.Percentage
	ValidatorEpochs       []model.ValidatorEpoch
	Accounts              []model.Account
	DelegatorBlockRewards []model.BlockReward
	Snarkers              []model.Snarker
	Transactions          []model.Transaction
	SnarkJobs             []model.SnarkJob
}
