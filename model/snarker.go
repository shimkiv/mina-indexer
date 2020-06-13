package model

import (
	"errors"
	"time"
)

type Snarker struct {
	Model

	Account    string `json:"account"`
	Fee        uint64 `json:"fee"`
	JobsCount  int    `json:"jobs_count"`
	WorksCount int    `json:"works_count"`

	StartHeight uint64    `json:"start_time"`
	StartTime   time.Time `json:"start_height"`
	LastHeight  uint64    `json:"last_height"`
	LastTime    time.Time `json:"last_time"`
}

func (s Snarker) Validate() error {
	if s.Account == "" {
		return errors.New("account is required")
	}
	return nil
}
