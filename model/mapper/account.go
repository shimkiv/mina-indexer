package mapper

import (
	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/model"
	"github.com/figment-networks/coda-indexer/model/util"
)

// Account returns an account model constructed from the coda input
func Account(block coda.Block, input coda.Account) (*model.Account, error) {
	startHeight, err := util.ParseInt64(block.ProtocolState.ConsensusState.BlockHeight)
	if err != nil {
		return nil, err
	}

	startedAt, err := util.ParseTime(block.ProtocolState.BlockchainState.Date)
	if err != nil {
		return nil, err
	}

	balance, err := util.ParseInt64(input.Balance.Total)
	if err != nil {
		return nil, err
	}

	var nonce int64
	if input.Nonce != nil {
		nonce, err = util.ParseInt64(*input.Nonce)
		if err != nil {
			return nil, err
		}
	}

	acc := &model.Account{
		PublicKey:   input.PublicKey,
		StartHeight: startHeight,
		StartedAt:   *startedAt,
		Balance:     balance,
		Nonce:       nonce,
	}

	return acc, acc.Validate()
}

// Accounts returns accounts models references from the block data
func Accounts(block coda.Block) ([]model.Account, error) {
	result := []model.Account{}

	validator, err := Account(block, *block.CreatorAccount)
	if err != nil {
		return nil, err
	}
	result = append(result, *validator)

	for _, t := range block.Transactions.UserCommands {
		from, err := Account(block, *t.FromAccount)
		if err != nil {
			return nil, err
		}
		to, err := Account(block, *t.ToAccount)
		if err != nil {
			return nil, err
		}
		result = append(result, *from, *to)
	}

	return result, nil
}
