package model

import (
	"errors"
	"fmt"
	"time"
)

// Block model contains block data
type Block struct {
	Model

	Time              time.Time `json:"time"`
	Height            int64     `json:"height"`
	Hash              string    `json:"hash"`
	ParentHash        string    `json:"parent_hash"`
	LedgerHash        string    `json:"ledger_hash"`
	Creator           string    `json:"creator"`
	TransactionsCount int       `json:"transactions_count"`
	Coinbase          int64     `json:"coinbase"`
	AppVersion        string    `json:"-"`
}

// BlockIntervalStat contains block count stats for a given time interval
type BlockIntervalStat struct {
	TimeInterval string  `json:"time_interval"`
	Count        int64   `json:"count"`
	Avg          float64 `json:"avg"`
}

// BlockAvgStat contains block averages
type BlockAvgStat struct {
	StartHeight int64   `json:"start_height"`
	EndHeight   int64   `json:"end_height"`
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	Count       int64   `json:"count"`
	Diff        float64 `json:"diff"`
	Avg         float64 `json:"avg"`
}

// TableName returns the model table name
func (Block) TableName() string {
	return "blocks"
}

// String returns a text representation of a block
func (b Block) String() string {
	return fmt.Sprintf("height=%v hash=%v", b.Height, b.Hash)
}

// Validate returns an error if block data is invalid
func (b Block) Validate() error {
	if b.Time.Year() == 0 {
		return errors.New("time is invalid")
	}
	if b.Height <= 0 {
		return errors.New("height is invalid")
	}
	if b.Hash == "" {
		return errors.New("hash is required")
	}
	if b.ParentHash == "" {
		return errors.New("parent hash is required")
	}
	if b.LedgerHash == "" {
		return errors.New("ledger hash is required")
	}
	if b.Creator == "" {
		return errors.New("creator is required")
	}
	return nil
}
