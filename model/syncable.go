package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	SyncableTypeBlock        = "block"
	SyncableTypeState        = "state"
	SyncableTypeTransactions = "transactions"
)

// Syncable contains raw blockchain data
type Syncable struct {
	Model

	ReportID    int64           `json:"report_id"`
	Height      int64           `json:"height"`
	Time        time.Time       `json:"time"`
	Type        string          `json:"type"`
	Data        json.RawMessage `json:"data"`
	ProcessedAt *time.Time      `json:"processed_at"`
}

// TableName returns the model table name
func (Syncable) TableName() string {
	return "syncables"
}

// String returns a text representation of syncable
func (s Syncable) String() string {
	return fmt.Sprintf("type=%v height=%v", s.Type, s.Height)
}

// Validate returns an error if syncable is invalid
func (s Syncable) Validate() error {
	if s.ReportID <= 0 {
		return errors.New("report id is required")
	}
	if s.Height <= 0 {
		return errors.New("height is invalid")
	}
	if s.Time.Year() == 0 {
		return errors.New("year is invalid")
	}
	if s.Type == "" {
		return errors.New("type is required")
	}
	if s.Data == nil {
		return errors.New("data is required")
	}
	return nil
}

// Decode decodes the raw data into a destination interface
func (s Syncable) Decode(dst interface{}) error {
	return json.Unmarshal(s.Data, dst)
}
