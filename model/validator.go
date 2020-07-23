package model

import (
	"errors"
	"time"
)

// Validator stores the block producer information
type Validator struct {
	Model

	Account        string `json:"account"`
	BlocksCreated  int    `json:"blocks_created"`
	BlocksProposed int    `json:"blocks_proposed"`

	StartHeight uint64    `json:"start_height"`
	StartTime   time.Time `json:"start_time"`
	LastHeight  uint64    `json:"last_height"`
	LastTime    time.Time `json:"last_time"`
}

type ValidatorStat struct {
	Time                string `json:"time"`
	Bucket              string `json:"bucket"`
	BlocksProducedCount int    `json:"blocks_produced_count"`
	DelegationsCount    int    `json:"delegations_count"`
	DelegationsAmount   int64  `json:"delegations_amount"`
}

// Validate returns an error if validator is invalid
func (v Validator) Validate() error {
	if v.Account == "" {
		return errors.New("account is required")
	}
	return nil
}
