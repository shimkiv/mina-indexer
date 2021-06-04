package indexing

import (
	"errors"

	"github.com/figment-networks/mina-indexer/client/archive"
	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/mapper"
	"github.com/figment-networks/mina-indexer/model/types"
)

// Prepare generates a new models from the graph block data
func Prepare(archiveBlock *archive.Block, graphBlock *graph.Block, validatorEpochs []model.ValidatorEpoch, ledgerData *mapper.LedgerData, firstSlotOfEpoch int) (*Data, error) {
	block, err := mapper.BlockFromArchive(archiveBlock)
	if err != nil {
		return nil, err
	}

	if graphBlock != nil {
		block.TotalCurrency = types.NewAmount(graphBlock.ProtocolState.ConsensusState.TotalCurrency)
		block.TransactionsFees = mapper.TransactionFees(graphBlock)
	}

	// Prepare validator record
	validator, err := mapper.Validator(archiveBlock)
	if err != nil {
		return nil, err
	}

	// Prepare transaction records
	transactions, err := mapper.TransactionsFromArchive(archiveBlock)
	if err != nil {
		return nil, err
	}
	block.TransactionsCount = len(transactions)

	// Prepare snarkers
	snarkers, err := mapper.Snarkers(graphBlock)
	if err != nil {
		return nil, err
	}
	block.SnarkersCount = len(snarkers)

	// Prepare snarker jobs
	snarkJobs, err := mapper.SnarkJobs(graphBlock)
	if err != nil {
		return nil, err
	}
	block.SnarkJobsCount = len(snarkJobs)
	block.SnarkJobsFees = types.NewInt64Amount(0)
	for _, job := range snarkJobs {
		block.SnarkJobsFees = block.SnarkJobsFees.Add(job.Fee)
	}

	// Prepare accounts
	accounts, err := mapper.Accounts(graphBlock)
	if err != nil {
		return nil, err
	}
	accountsMap := map[string]*model.Account{}
	for _, acc := range accounts {
		accountsMap[acc.PublicKey] = &acc
	}
	var validatorBlockReward *model.BlockReward
	delegatorBlockRewards := []model.BlockReward{}
	var creatorAcc *model.Account
	var ok bool
	var supercharged bool
	if graphBlock != nil {
		creatorAcc, ok = accountsMap[graphBlock.Creator]
		if !ok {
			return nil, errors.New("creator is not found in accounts map " + graphBlock.Creator)
		}
		validatorBlockReward, _ = mapper.ValidatorBlockReward(validator)
		delegatorBlockRewards, err = mapper.DelegatorBlockRewards(ledgerData.Entries, graphBlock)
		if err != nil {
			return nil, err
		}

		// the reward is supercharged is based off of the account that won the block
		if graphBlock.WinnerAccount.Locked != nil {
			supercharged = !(*graphBlock.WinnerAccount.Locked)
		} else {
			// The coinbase reward for producing a block is 720 tokens. for supercharged 2x
			supercharged = block.Coinbase.Int64() == 1440000000000
		}
	}

	data := &Data{
		Block:                  block,
		FirstSlotOfEpoch:       firstSlotOfEpoch,
		Supercharged:           supercharged,
		Validator:              validator,
		ValidatorBlockReward:   validatorBlockReward,
		CreatorAccount:         creatorAcc,
		ValidatorEpochs:        validatorEpochs,
		Accounts:               accounts,
		DelegatorsBlockRewards: delegatorBlockRewards,
		Transactions:           transactions,
		Snarkers:               snarkers,
		SnarkJobs:              snarkJobs,
	}

	return data, nil
}
