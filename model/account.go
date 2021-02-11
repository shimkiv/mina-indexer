package model

import (
	"errors"
	"time"

	"github.com/figment-networks/mina-indexer/model/types"
)

// Account contains the account details
type Account struct {
	ID             string       `json:"-"`
	PublicKey      string       `json:"public_key"`
	Delegate       *string      `json:"delegate"`
	Balance        types.Amount `json:"balance"`
	BalanceUnknown types.Amount `json:"balance_unknown"`
	Stake          types.Amount `json:"stake"`
	Nonce          uint64       `json:"nonce"`
	StartHeight    uint64       `json:"start_height"`
	StartTime      time.Time    `json:"start_time"`
	LastHeight     uint64       `json:"last_height"`
	LastTime       time.Time    `json:"last_time"`
	CreatedAt      time.Time    `json:"-"`
	UpdatedAt      time.Time    `json:"-"`
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
	return nil
}
