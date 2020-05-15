package mapper

import (
	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/model"
	"github.com/figment-networks/coda-indexer/model/util"
)

// Transaction returns a transaction model constructed from the coda input
func Transaction(block coda.Block, t coda.UserCommand) (*model.Transaction, error) {
	ttype := model.TransactionTypePayment
	if t.IsDelegation {
		ttype = model.TransactionTypeDelegation
	}

	tran := &model.Transaction{
		Type:         ttype,
		Time:         blockTime(block),
		Height:       blockHeight(block),
		Hash:         t.ID,
		BlockHash:    block.StateHash,
		SenderKey:    t.From,
		RecipientKey: t.To,
		Amount:       util.MustInt64(t.Amount),
		Fee:          util.MustInt64(t.Fee),
		Nonce:        int64(t.Nonce),
	}

	return tran, tran.Validate()
}

// Transactions returns a list of transactions from the coda input
func Transactions(block coda.Block) ([]model.Transaction, error) {
	if block.Transactions == nil {
		return nil, nil
	}
	if block.Transactions.UserCommands == nil {
		return nil, nil
	}

	commands := block.Transactions.UserCommands
	result := make([]model.Transaction, len(commands))

	for i, cmd := range commands {
		t, err := Transaction(block, *cmd)
		if err != nil {
			return nil, err
		}
		result[i] = *t
	}

	return result, nil
}
