package model

import (
	"errors"
	"time"
)

// Account contains the account details
type Account struct {
	Model

	// Name represents a reference to the account owner, it does not exist
	// on the chain and should be imported from the annotated ledger.
	// Name *string `json:"name"`

	PublicKey      string  `json:"public_key"`
	Delegate       *string `json:"delegate"`
	Balance        string  `json:"balance"`
	BalanceUnknown string  `json:"balance_unknown"`
	Nonce          int64   `json:"nonce"`
	//VotingFor      string  `json:"voting_for"`
	//TxSent         int     `json:"tx_sent"`
	//TxReceived     int     `json:"tx_received"`

	StartHeight uint64    `json:"start_time"`
	StartTime   time.Time `json:"start_height"`
	LastHeight  uint64    `json:"last_height"`
	LastTime    time.Time `json:"last_time"`
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
	if acc.Balance == "" {
		return errors.New("balance is required")
	}
	return nil
}
