package model

import (
	"errors"
	"fmt"
	"time"
)

// Job contains a completed SNARK job details
type Job struct {
	Model

	Height     uint64    `json:"height"`
	Time       time.Time `json:"time"`
	Prover     string    `json:"prover"`
	Fee        uint64    `json:"fee"`
	WorksCount int       `json:"works_count"`
}

// TableName returns the Job table name
func (Job) TableName() string {
	return "jobs"
}

func (j Job) String() string {
	return fmt.Sprintf(
		"height=%v prover=%s fee=%v",
		j.Height,
		j.Prover,
		j.Fee,
	)
}

// Validate returns an error if job is invalid
func (j Job) Validate() error {
	if j.Height <= 0 {
		return errors.New("height is invalid")
	}
	if j.Time.Year() == 0 {
		return errors.New("time is invalid")
	}
	if j.Prover == "" {
		return errors.New("prover is required")
	}
	if j.Fee < 0 {
		return errors.New("fee is invalid")
	}
	return nil
}
