package model

import (
	"errors"
	"time"

	"github.com/figment-networks/mina-indexer/model/types"
)

// Validator stores the block producer information
type Validator struct {
	ID                int          `json:"-"`
	IdentityName      *string      `json:"identity_name"`
	PublicKey         string       `json:"public_key"`
	BlocksCreated     int          `json:"blocks_created"`
	BlocksProposed    int          `json:"blocks_proposed"`
	DelegatedAccounts int          `json:"delegated_accounts"`
	DelegatedBalance  types.Amount `json:"delegated_balance"`
	Stake             types.Amount `json:"stake"`
	StartHeight       uint64       `json:"start_height"`
	StartTime         time.Time    `json:"start_time"`
	LastHeight        uint64       `json:"last_height"`
	LastTime          time.Time    `json:"last_time"`
	CreatedAt         time.Time    `json:"-"`
	UpdatedAt         time.Time    `json:"-"`
}

type ValidatorStat struct {
	Time                string `json:"time"`
	Bucket              string `json:"bucket"`
	BlocksProducedCount int    `json:"blocks_produced_count"`
	DelegationsCount    int    `json:"delegations_count"`
	DelegationsAmount   string `json:"delegations_amount"`
}

// Validate returns an error if validator is invalid
func (v Validator) Validate() error {
	if v.PublicKey == "" {
		return errors.New("account is required")
	}
	return nil
}
