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

	if input.Delegate != nil && *input.Delegate != input.PublicKey {
		acc.Delegate = input.Delegate
	}

	if input.Nonce != nil {
		acc.Nonce = util.MustInt64(*input.Nonce)
	}

	return acc, acc.Validate()
}

// Accounts returns accounts models references from the block data
func Accounts(block *coda.Block) ([]model.Account, error) {
	codaAccounts := []*coda.Account{
		block.CreatorAccount,
		block.CreatorAccount.DelegateAccount,
	}
	for _, t := range block.Transactions.UserCommands {
		codaAccounts = append(codaAccounts,
			t.FromAccount, t.FromAccount.DelegateAccount,
			t.ToAccount, t.ToAccount.DelegateAccount,
		)
	}

	accounts := map[string]*model.Account{}
	for _, codaAcc := range codaAccounts {
		if codaAcc == nil {
			continue
		}

		acc, err := Account(block, codaAcc)
		if err != nil {
			return nil, err
		}

		if accounts[acc.PublicKey] == nil {
			accounts[acc.PublicKey] = acc
		}
	}

	result := []model.Account{}
	for _, v := range accounts {
		result = append(result, *v)
	}

	return result, nil
}
