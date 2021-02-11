package model

import (
	"errors"
	"time"

	"github.com/figment-networks/mina-indexer/model/types"
)

const (
	// Transaction types
	TxTypePayment             = "payment"
	TxTypeDelegation          = "delegation"
	TxTypeCoinbase            = "coinbase"
	TxTypeCoinbaseFeeTransfer = "fee_transfer_via_coinbase"
	TxTypeFeeTransfer         = "fee_transfer"
	TxTypeSnarkFee            = "snark_fee"

	// Transaction statuses
	TxStatusApplied = "applied"
	TxStatusFailed  = "failed"
)

var (
	TxTypes = []string{
		TxTypePayment,
		TxTypeDelegation,
		TxTypeCoinbase,
		TxTypeCoinbaseFeeTransfer,
		TxTypeFeeTransfer,
		TxTypeSnarkFee,
	}
)

// Transaction contains the blockchain transaction details
type Transaction struct {
	ID                      int          `json:"id"`
	Hash                    string       `json:"hash"`
	Type                    string       `json:"type"`
	BlockHash               string       `json:"block_hash"`
	BlockHeight             uint64       `json:"block_height"`
	Time                    time.Time    `json:"time"`
	Sender                  *string      `json:"sender"`
	Receiver                string       `json:"receiver"`
	Amount                  types.Amount `json:"amount"`
	Fee                     types.Amount `json:"fee"`
	Nonce                   *int         `json:"nonce"`
	Memo                    *string      `json:"memo"`
	Status                  string       `json:"status"`
	FailureReason           *string      `json:"failure_reason"`
	SequenceNumber          *int         `json:"sequence_number"`
	SecondarySequenceNumber *int         `json:"secondary_sequence_number"`
	CreatedAt               time.Time    `json:"-"`
	UpdatedAt               time.Time    `json:"-"`
}

// TableName returns the model table name
func (Transaction) TableName() string {
	return "transactions"
}

// Validate returns an error if transaction is invalid
func (t Transaction) Validate() error {
	if t.Type == "" {
		return errors.New("type is required")
	}
	if t.BlockHash == "" {
		return errors.New("block hash is required")
	}
	if t.Hash == "" {
		return errors.New("hash is required")
	}
	if t.BlockHeight <= 0 {
		return errors.New("height is invalid")
	}
	if t.Time.IsZero() {
		return errors.New("time is invalid")
	}
	if t.Receiver == "" {
		return errors.New("receiver is required")
	}
	return nil
}
