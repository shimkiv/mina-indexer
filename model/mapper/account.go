package mapper

import (
	"github.com/figment-networks/mina-indexer/client/archive"
	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/types"
	"github.com/figment-networks/mina-indexer/model/util"
)

// Account returns an account model constructed from the graph input
func Account(block *graph.Block, input *graph.Account) (*model.Account, error) {
	height := BlockHeight(block)
	time := BlockTime(block)

	acc := &model.Account{
		PublicKey:      input.PublicKey,
		StartHeight:    height,
		StartTime:      time,
		LastHeight:     height,
		LastTime:       time,
		Balance:        types.NewAmount(input.Balance.Total),
		BalanceUnknown: types.NewAmount(input.Balance.Unknown),
	}

	if input.Delegate != nil && *input.Delegate != input.PublicKey {
		acc.Delegate = input.Delegate
	}

	if input.Nonce != nil {
		acc.Nonce = util.MustUInt64(*input.Nonce)
	}

	if input.Locked != nil {
		acc.Supercharged = !(*input.Locked)
	}

	return acc, acc.Validate()
}

// Accounts returns accounts models references from the block data
func Accounts(block *graph.Block) ([]model.Account, error) {
	if block == nil {
		return nil, nil
	}

	graphAccounts := []*graph.Account{
		block.CreatorAccount,
		block.CreatorAccount.DelegateAccount,
	}

	accounts := map[string]*model.Account{}
	for _, graphAcc := range graphAccounts {
		if graphAcc == nil {
			continue
		}

		acc, err := Account(block, graphAcc)
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

func AccountFromStagedLedger(block *graph.Block, entry *archive.StakingInfo) (*model.Account, error) {
	height := BlockHeight(block)
	time := BlockTime(block)

	account := &model.Account{
		PublicKey:      entry.Pk,
		Balance:        types.NewFloatAmount(entry.Balance),
		BalanceUnknown: types.NewFloatAmount(entry.Balance),
		StartHeight:    height,
		StartTime:      time,
		LastHeight:     height,
		LastTime:       time,
	}

	if entry.Pk != entry.Delegate {
		delegate := entry.Delegate
		account.Delegate = &delegate
	} else {
		account.Delegate = nil
	}

	return account, nil
}
