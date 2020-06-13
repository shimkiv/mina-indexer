package mapper

import (
	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/model"
	"github.com/figment-networks/coda-indexer/model/util"
)

// Account returns an account model constructed from the coda input
func Account(block *coda.Block, input *coda.Account) (*model.Account, error) {
	height := BlockHeight(block)
	time := BlockTime(block)

	acc := &model.Account{
		PublicKey:      input.PublicKey,
		StartHeight:    height,
		StartTime:      time,
		LastHeight:     height,
		LastTime:       time,
		Balance:        input.Balance.Total,
		BalanceUnknown: input.Balance.Unknown,
	}

	if input.Nonce != nil {
		acc.Nonce = util.MustInt64(*input.Nonce)
	}

	return acc, acc.Validate()
}

// Accounts returns accounts models references from the block data
func Accounts(block *coda.Block) ([]model.Account, error) {
	accounts := map[string]*model.Account{}

	// Prepare validator record
	validator, err := Account(block, block.CreatorAccount)
	if err != nil {
		return nil, err
	}
	accounts[validator.PublicKey] = validator

	// Prepare accounts from user transactions
	for _, t := range block.Transactions.UserCommands {
		from, err := Account(block, t.FromAccount)
		if err != nil {
			return nil, err
		}
		if accounts[from.PublicKey] == nil {
			accounts[from.PublicKey] = from
		}

		to, err := Account(block, t.ToAccount)
		if err != nil {
			return nil, err
		}
		if accounts[to.PublicKey] == nil {
			accounts[to.PublicKey] = to
		}
	}

	result := []model.Account{}
	for _, v := range accounts {
		result = append(result, *v)
	}

	return result, nil
}
