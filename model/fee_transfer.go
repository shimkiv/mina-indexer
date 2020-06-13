package model

import (
	"errors"
	"time"
)

type FeeTransfer struct {
	Model

	Height    uint64    `json:"height"`
	Time      time.Time `json:"time"`
	Recipient string    `json:"recipient"`
	Amount    uint64    `json:"amount"`
}

func (t FeeTransfer) Validate() error {
	if t.Height == 0 {
		return errors.New("height is invalid")
	}
	if t.Time.IsZero() {
		return errors.New("time is invalid")
	}
	if t.Recipient == "" {
		return errors.New("recipient is required")
	}
	return nil
}
