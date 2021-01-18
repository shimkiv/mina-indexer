package indexing

import (
	"github.com/figment-networks/mina-indexer/coda"
	"github.com/figment-networks/mina-indexer/model/mapper"
)

// Prepare generates a new models from the coda block data
func Prepare(status *coda.DaemonStatus, input *coda.Block) (*Data, error) {
	// Prepare block record
	block, err := mapper.Block(input)
	if err != nil {
		return nil, err
	}
	block.AppVersion = status.CommitID

	// Prepare validator record
	validator, err := mapper.Validator(input)
	if err != nil {
		return nil, err
	}

	// Prepare transaction records
	transactions, err := mapper.Transactions(input)
	if err != nil {
		return nil, err
	}

	// Prepare fee transfers
	transfers, err := mapper.FeeTransfers(input)
	if err != nil {
		return nil, err
	}

	// Prepare snarkers
	snarkers, err := mapper.Snarkers(input)
	if err != nil {
		return nil, err
	}

	// Prepare snarker jobs
	jobs, err := mapper.Jobs(input)
	if err != nil {
		return nil, err
	}

	// Prepare accounts
	accounts, err := mapper.Accounts(input)
	if err != nil {
		return nil, err
	}

	data := &Data{
		Block:        block,
		Validator:    validator,
		Accounts:     accounts,
		Transactions: transactions,
		FeeTransfers: transfers,
		Snarkers:     snarkers,
		SnarkJobs:    jobs,
	}

	return data, nil
}
