package model

import (
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

	Hash         string    `json:"hash"`
	Type         string    `json:"type"`
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
