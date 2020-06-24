package mapper

import (
	"fmt"

	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/model"
	"github.com/figment-networks/coda-indexer/model/util"
)

func UserTransaction(block *coda.Block, t *coda.UserCommand) (*model.Transaction, error) {
	ttype := model.TxTypePayment
	if t.IsDelegation {
		ttype = model.TxTypeDelegation
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

func BlockRewardTransaction(block *coda.Block) (*model.Transaction, error) {
	t := &model.Transaction{
		Type:      model.TxTypeBlockReward,
		BlockHash: block.StateHash,
		Hash:      util.SHA1(block.StateHash + block.Transactions.CoinbaseReceiver.PublicKey),
		Height:    BlockHeight(block),
		Time:      BlockTime(block),
		Receiver:  block.Transactions.CoinbaseReceiver.PublicKey,
		Amount:    util.MustUInt64(block.Transactions.Coinbase),
	}

	return t, t.Validate()
}

func SnarkFeeTransaction(block *coda.Block, transfer *coda.FeeTransfer) (*model.Transaction, error) {
	uid := fmt.Sprintf("%s%s%s", block.StateHash, transfer.Recipient, transfer.Fee)

	t := &model.Transaction{
		Type:      model.TxTypeSnarkFee,
		BlockHash: block.StateHash,
		Hash:      util.SHA1(uid),
		Height:    BlockHeight(block),
		Time:      BlockTime(block),
		Receiver:  transfer.Recipient,
		Amount:    util.MustUInt64(transfer.Fee),
	}

	return t, t.Validate()
}

// Transactions returns a list of transactions from the coda input
func Transactions(block *coda.Block) ([]model.Transaction, error) {
	if block.Transactions == nil {
		return nil, nil
	}
	if block.Transactions.UserCommands == nil {
		return nil, nil
	}

	result := []model.Transaction{}

	// Add the block reward transaction
	if t, err := BlockRewardTransaction(block); err == nil {
		result = append(result, *t)
	} else {
		return nil, err
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
	feeTransfers := block.Transactions.FeeTransfer
	for _, transfer := range feeTransfers {
		t, err := SnarkFeeTransaction(block, transfer)
		if err != nil {
			return nil, err
		}
		result = append(result, *t)
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
