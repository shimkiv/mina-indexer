package model

import (
	"errors"
	"fmt"
	"time"
)

const (
	TransactionTypePayment    = "payment"
	TransactionTypeDelegation = "delegation"
)

// Transaction contains the blockchain transaction details
type Transaction struct {
	Model

	Type         string    `json:"type"`
	Hash         string    `json:"hash"`
	BlockHash    string    `json:"block_hash"`
	Height       int64     `json:"height"`
	Time         time.Time `json:"time"`
	Nonce        int64     `json:"nonce"`
	SenderKey    string    `json:"sender_key"`
	RecipientKey string    `json:"recipient_key"`
	Amount       int64     `json:"amount"`
	Fee          int64     `json:"fee"`
}

// TableName returns the model table name
func (Transaction) TableName() string {
	return "transactions"
}

// String returns transaction text representation
func (t Transaction) String() string {
	return fmt.Sprintf("type=%v hash=%v height=%v", t.Type, t.Hash, t.Height)
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
	if t.Height <= 0 {
		return errors.New("height is invalid")
	}
	if t.Time.Year() == 0 {
		return errors.New("time is invalid")
	}
	if t.SenderKey == "" {
		return errors.New("sender key is required")
	}
	if t.RecipientKey == "" {
		return errors.New("recipient key is required")
	}
	if t.Amount < 0 {
		return errors.New("amount is invalid")
	}
	return nil
}
