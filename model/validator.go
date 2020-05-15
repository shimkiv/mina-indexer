package model

import (
	"errors"
	"time"
)

// Validator stores the block producer information
type Validator struct {
	Model

	PublicKey string    `json:"public_key"`
	Height    int64     `json:"height"`
	Time      time.Time `json:"time"`
}

// Validate returns an error if validator is invalid
func (v Validator) Validate() error {
	if v.PublicKey == "" {
		return errors.New("public key is required")
	}
	if v.Height <= 0 {
		return errors.New("height is invalid")
	}
	if v.Time.Year() == 0 {
		return errors.New("time is invalid")
	}
	return nil
}
