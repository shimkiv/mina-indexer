package mapper

import (
	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/model"
	"github.com/figment-networks/coda-indexer/model/util"
)

// Transaction returns a transaction model constructed from the coda input
func Transaction(block *coda.Block, t *coda.UserCommand) (*model.Transaction, error) {
	ttype := model.TransactionTypePayment
	if t.IsDelegation {
		ttype = model.TransactionTypeDelegation
	}

	tran := &model.Transaction{
		Type:      ttype,
		Time:      BlockTime(block),
		Height:    BlockHeight(block),
		Hash:      t.ID,
		BlockHash: block.StateHash,
		Sender:    t.From,
		Receiver:  t.To,
		Amount:    util.MustUInt64(t.Amount),
		Fee:       util.MustUInt64(t.Fee),
		Nonce:     t.Nonce,
		Memo:      t.Memo,
	}

	return tran, tran.Validate()
}

// Transactions returns a list of transactions from the coda input
func Transactions(block *coda.Block) ([]model.Transaction, error) {
	if block.Transactions == nil {
		return nil, nil
	}
	if block.Transactions.UserCommands == nil {
		return nil, nil
	}

	commands := block.Transactions.UserCommands
	result := make([]model.Transaction, len(commands))

	for i, cmd := range commands {
		t, err := Transaction(block, cmd)
		if err != nil {
			return nil, err
		}
		result[i] = *t
	}

	return result, nil
}

func FeeTransfers(block *coda.Block) ([]model.FeeTransfer, error) {
	if block.Transactions == nil {
		return nil, nil
	}
	if block.Transactions.FeeTransfer == nil {
		return nil, nil
	}

	transfers := block.Transactions.FeeTransfer
	result := make([]model.FeeTransfer, len(transfers))

	for i, t := range transfers {
		result[i].Height = BlockHeight(block)
		result[i].Time = BlockTime(block)
		result[i].Recipient = t.Recipient
		result[i].Amount = util.MustUInt64(t.Fee)
	}

	return result, nil
}
