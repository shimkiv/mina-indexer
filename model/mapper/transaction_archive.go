package mapper

import (
	"time"

	"github.com/figment-networks/mina-indexer/client/archive"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/types"
	"github.com/figment-networks/mina-indexer/model/util"
)

func TransactionsFromArchive(block *archive.Block) ([]model.Transaction, error) {
	blockHeight := uint64(block.Height)
	blockTime := time.Unix(block.Timestamp/1000, 0)
	result := make([]model.Transaction, len(block.UserCommands)+len(block.InternalCommands))
	idx := 0

	for _, cmd := range block.InternalCommands {
		result[idx] = model.Transaction{
			Type:                    cmd.Type,
			Hash:                    cmd.ID,
			BlockHash:               block.StateHash,
			BlockHeight:             blockHeight,
			Time:                    blockTime,
			Receiver:                cmd.Receiver,
			Amount:                  types.NewInt64Amount(cmd.Fee),
			Status:                  model.TxStatusApplied,
			SequenceNumber:          &cmd.SequenceNo,
			SecondarySequenceNumber: &cmd.SecondarySequenceNo,
		}
		idx++
	}

	for _, cmd := range block.UserCommands {
		var memoText *string
		if text := util.ParseMemoText(cmd.Memo); len(text) > 0 {
			memoText = &text
		}

		result[idx] = model.Transaction{
			Type:           cmd.Type,
			Hash:           cmd.Hash,
			BlockHash:      block.StateHash,
			BlockHeight:    blockHeight,
			Time:           blockTime,
			Sender:         &cmd.Sender,
			Receiver:       cmd.Receiver,
			Amount:         types.NewInt64Amount(cmd.Amount),
			Fee:            types.NewInt64Amount(cmd.Fee),
			Status:         cmd.Status,
			FailureReason:  cmd.FailureReason,
			SequenceNumber: &cmd.SequenceNo,
			Nonce:          &cmd.Nonce,
			Memo:           memoText,
		}
		idx++
	}

	return result, nil
}
