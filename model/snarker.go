package model

import (
	"errors"
	"time"

	"github.com/figment-networks/mina-indexer/model/types"
)

type Snarker struct {
	ID          int       `json:"-"`
	Account     string    `json:"public_key"`
	Fee         uint64    `json:"fee"`
	JobsCount   int       `json:"jobs_count"`
	WorksCount  int       `json:"works_count"`
	StartHeight uint64    `json:"start_height"`
	StartTime   time.Time `json:"start_time"`
	LastHeight  uint64    `json:"last_height"`
	LastTime    time.Time `json:"last_time"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

func (s Snarker) Validate() error {
	if s.Account == "" {
		return errors.New("public key is required")
	}
	return nil
}

type SnarkerJobsStat struct {
	Time          time.Time    `json:"time"`
	FeeMin        types.Amount `json:"fee_min"`
	FeeMax        types.Amount `json:"fee_max"`
	FeeAvg        types.Amount `json:"fee_avg"`
	JobsCount     int          `json:"jobs_count"`
	SnarkersCount int          `json:"snarkers_count"`
	WorksCount    int          `json:"works_count"`
}
