package model

import (
	"errors"
	"fmt"
)

// State stores the consensus state attributes
type State struct {
	Model

	Height        int64  `json:"height"`
	TotalCurrency int64  `json:"total_currency"`
	Epoch         int64  `json:"epoch"`
	EpochCount    int64  `json:"epoch_count"`
	LastVFROutput string `json:"last_vfr_output"`
}

// String returns state text representation
func (s State) String() string {
	return fmt.Sprintf("height=%v currency=%v", s.Height, s.TotalCurrency)
}

// Validate validates the state record
func (s State) Validate() error {
	if s.Height <= 0 {
		return errors.New("height is invalid")
	}
	if s.TotalCurrency <= 0 {
		return errors.New("total currency is invalid")
	}
	return nil
}
