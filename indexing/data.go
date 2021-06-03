package indexing

import (
	"github.com/figment-networks/mina-indexer/model"
)

// Data contains all the records processed for a height
type Data struct {
	Block                  *model.Block
	FirstSlotOfEpoch       int
	Supercharged           bool
	Validator              *model.Validator
	ValidatorBlockReward   *model.BlockReward
	CreatorAccount         *model.Account
	ValidatorEpochs        []model.ValidatorEpoch
	Accounts               []model.Account
	DelegatorsBlockRewards []model.BlockReward
	Snarkers               []model.Snarker
	Transactions           []model.Transaction
	SnarkJobs              []model.SnarkJob
}
