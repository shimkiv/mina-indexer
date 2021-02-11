package model

import (
	"errors"
	"time"
)

type Snarker struct {
	ID          int       `json:"-"`
	PublicKey   string    `json:"public_key"`
	Fee         uint64    `json:"fee"`
	JobsCount   int       `json:"jobs_count"`
	WorksCount  int       `json:"works_count"`
	StartHeight uint64    `json:"start_time"`
	StartTime   time.Time `json:"start_height"`
	LastHeight  uint64    `json:"last_height"`
	LastTime    time.Time `json:"last_time"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

func (s Snarker) Validate() error {
	if s.PublicKey == "" {
		return errors.New("public key is required")
	}
	return nil
}
