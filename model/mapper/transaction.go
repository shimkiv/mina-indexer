package mapper

import (
	"fmt"

	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/types"
	"github.com/figment-networks/mina-indexer/model/util"
)

func UserTransaction(block *graph.Block, t *graph.UserCommand) (*model.Transaction, error) {
	ttype := model.TxTypePayment
	if t.IsDelegation {
		ttype = model.TxTypeDelegation
	}

	var memoText *string
	if text := util.ParseMemoText(t.Memo); len(text) > 0 {
		memoText = &text
	}

	tran := &model.Transaction{
		Type:        ttype,
		Hash:        t.Hash,
		Time:        BlockTime(block),
		BlockHeight: BlockHeight(block),
		BlockHash:   block.StateHash,
		Sender:      &t.From,
		Receiver:    t.To,
		Amount:      types.NewAmount(t.Amount),
		Fee:         types.NewAmount(t.Fee),
		Nonce:       &t.Nonce,
		Memo:        memoText,
	}

	return tran, tran.Validate()
}

func BlockRewardTransaction(block *graph.Block) (*model.Transaction, error) {
	t := &model.Transaction{
		Type:        model.TxTypeCoinbase,
		Hash:        util.SHA1(block.StateHash + block.Transactions.CoinbaseReceiver.PublicKey),
		Time:        BlockTime(block),
		BlockHash:   block.StateHash,
		BlockHeight: BlockHeight(block),
		Receiver:    block.Transactions.CoinbaseReceiver.PublicKey,
		Amount:      types.NewAmount(block.Transactions.Coinbase),
		Fee:         types.NewInt64Amount(0),
	}

	return t, t.Validate()
}

func FeeTransaction(block *graph.Block, transfer *graph.FeeTransfer) (*model.Transaction, error) {
	uid := fmt.Sprintf("%s%s%s", block.StateHash, transfer.Recipient, transfer.Fee)

	t := &model.Transaction{
		Type:        model.TxTypeFeeTransfer,
		Hash:        util.SHA1(uid),
		Time:        BlockTime(block),
		BlockHash:   block.StateHash,
		BlockHeight: BlockHeight(block),
		Receiver:    transfer.Recipient,
		Amount:      types.NewAmount(transfer.Fee),
		Fee:         types.NewInt64Amount(0),
	}

	return t, t.Validate()
}

func SnarkFeeTransaction(block *graph.Block, transfer *graph.FeeTransfer) (*model.Transaction, error) {
	uid := fmt.Sprintf("%s%s%s", block.StateHash, transfer.Recipient, transfer.Fee)

	t := &model.Transaction{
		Type:        model.TxTypeSnarkFee,
		Hash:        util.SHA1(uid),
		Time:        BlockTime(block),
		BlockHash:   block.StateHash,
		BlockHeight: BlockHeight(block),
		Sender:      &block.Creator,
		Receiver:    transfer.Recipient,
		Amount:      types.NewAmount(transfer.Fee),
		Fee:         types.NewInt64Amount(0),
	}

	return t, t.Validate()
}

// Transactions returns a list of transactions from the coda input
func Transactions(block *graph.Block) ([]model.Transaction, error) {
	if block.Transactions == nil {
		return nil, nil
	}
	if block.Transactions.UserCommands == nil {
		return nil, nil
	}

	result := []model.Transaction{}

	// Add the block reward transaction
	if block.Transactions.Coinbase != "0" && block.Transactions.CoinbaseReceiver != nil {
		if t, err := BlockRewardTransaction(block); err == nil {
			result = append(result, *t)
		} else {
			return nil, err
		}
	}

	// Add user transactions
	commands := block.Transactions.UserCommands
	for _, cmd := range commands {
		t, err := UserTransaction(block, cmd)
		if err != nil {
			return nil, err
		}
		result = append(result, *t)
	}

	// Add snarker fees transactions
	snarkerIDs := map[string]bool{}
	for _, job := range block.SnarkJobs {
		snarkerIDs[job.Prover] = true
	}

	feeTransfers := block.Transactions.FeeTransfer
	for _, transfer := range feeTransfers {
		var feeTx *model.Transaction
		var err error

		if snarkerIDs[transfer.Recipient] {
			feeTx, err = SnarkFeeTransaction(block, transfer)
		} else {
			feeTx, err = FeeTransaction(block, transfer)
		}
		if err != nil {
			return nil, err
		}

		result = append(result, *feeTx)
	}

	return result, nil
}
