package model

import (
	"errors"
	"time"

	"github.com/figment-networks/mina-indexer/model/types"
)

// SnarkJob contains a completed SNARK job details
type SnarkJob struct {
	ID             int          `json:"-"`
	Height         uint64       `json:"height"`
	BlockReference string       `json:"block_reference"`
	Time           time.Time    `json:"time"`
	Prover         string       `json:"prover"`
	Fee            types.Amount `json:"fee"`
	WorksCount     int          `json:"works_count"`
	CreatedAt      time.Time    `json:"-"`
}

// TableName returns the Job table name
func (SnarkJob) TableName() string {
	return "snark_jobs"
}

// Validate returns an error if job is invalid
func (j SnarkJob) Validate() error {
	if j.Height <= 0 {
		return errors.New("height is invalid")
	}
	if j.Time.IsZero() {
		return errors.New("time is invalid")
	}
	if j.Prover == "" {
		return errors.New("prover is required")
	}
	return nil
}
