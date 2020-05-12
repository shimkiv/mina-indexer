package model

import (
	"errors"
	"time"
)

// Account contains the account details
type Account struct {
	Model

	PublicKey   string    `json:"public_key"`
	StartHeight int64     `json:"start_height"`
	StartedAt   time.Time `json:"started_at"`
	Balance     int64     `json:"balance"`
	Nonce       int64     `json:"nonce"`
}

// TableName returns the model table name
func (Account) TableName() string {
	return "accounts"
}

// String returns account text representation
func (acc Account) String() string {
	return acc.PublicKey
}

// Validate validates the record
func (acc Account) Validate() error {
	if acc.PublicKey == "" {
		return errors.New("public key is required")
	}
	if acc.StartHeight <= 0 {
		return errors.New("start height is invalid")
	}
	if acc.StartedAt.Year() == 0 {
		return errors.New("start date is invalid")
	}
	if acc.Balance < 0 {
		return errors.New("balance is invalid")
	}
	return nil
}
