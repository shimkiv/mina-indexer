package indexing

import (
	"github.com/figment-networks/mina-indexer/model"
)

// Data contains all the records processed for a height
type Data struct {
	Block           *model.Block
	Validator       *model.Validator
	CreatorAccount  *model.Account
	ValidatorEpochs []model.ValidatorEpoch
	Accounts        []model.Account
	Snarkers        []model.Snarker
	Transactions    []model.Transaction
	SnarkJobs       []model.SnarkJob
}

// AccountIDs returns a list of all accounts seen in the data payload
func (d Data) AccountIDs() []string {
	accounts := map[string]bool{
		d.Block.Creator: true,
	}

	for _, tx := range d.Transactions {
		accounts[tx.Receiver] = true
		if tx.Sender != nil {
			accounts[*tx.Sender] = true
		}
	}

	for _, snarker := range d.Snarkers {
		accounts[snarker.Account] = true
	}

	result := make([]string, len(accounts))
	idx := 0

	for k := range accounts {
		result[idx] = k
		idx++
	}

	return result
}
