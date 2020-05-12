package model

import (
	"encoding/json"
	"errors"
	"time"
)

const (
	ReportStatePending  = "pending"
	ReportStateRunning  = "running"
	ReportStateFinished = "finished"
)

// Report contains the sync run details
type Report struct {
	*Model

	StartHeight  int64
	EndHeight    int64
	State        string
	SuccessCount *int64
	ErrorCount   *int64
	ErrorMsg     *string
	Duration     *int64
	Details      *json.RawMessage
	CompletedAt  *time.Time
}

// TableName returns the model table name
func (Report) TableName() string {
	return "reports"
}

// Validate returns an error if report is invalid
func (r *Report) Validate() error {
	if r.StartHeight <= 0 {
		return errors.New("start height is invalid")
	}
	if r.EndHeight <= 0 {
		return errors.New("end height is invalid")
	}
	if r.StartHeight >= r.EndHeight {
		return errors.New("start height is greater that end height")
	}
	if r.State == "" {
		return errors.New("state is required")
	}
	return nil
}

// Complete marks the report as completed
func (r *Report) Complete(successCount, errCount int64, err *string, details []byte) {
	now := time.Now()
	duration := now.Sub(r.CreatedAt).Milliseconds()

	r.State = ReportStateFinished
	r.CompletedAt = &now
	r.Duration = &duration
	r.SuccessCount = &successCount
	r.ErrorCount = &errCount
}
