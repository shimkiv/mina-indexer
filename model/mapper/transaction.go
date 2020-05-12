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

	time, err := util.ParseTime(block.ProtocolState.BlockchainState.Date)
	if err != nil {
		return nil, err
	}

	height, err := util.ParseInt64(block.ProtocolState.ConsensusState.BlockHeight)
	if err != nil {
		return nil, err
	}

	amount, err := util.ParseInt64(t.Amount)
	if err != nil {
		return nil, err
	}

	fee, err := util.ParseInt64(t.Fee)
	if err != nil {
		return nil, err
	}

	tran := &model.Transaction{
		Hash:         t.ID,
		Type:         ttype,
		Time:         *time,
		Height:       height,
		SenderKey:    t.From,
		RecipientKey: t.To,
		Amount:       amount,
		Fee:          fee,
		Nonce:        int64(t.Nonce),
	}

	return tran, nil
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
