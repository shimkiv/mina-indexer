package store

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/figment-networks/coda-indexer/model"
)

var (
	reDate = regexp.MustCompile(`^[\d]{4}-[\d]{2}-[\d]{2}$`)
)

// TransactionSearch contains transaction search params
type TransactionSearch struct {
	AfterID   uint   `form:"after_id"`
	BeforeID  uint   `form:"before_id"`
	Height    uint64 `form:"height"`
	Type      string `form:"type"`
	BlockHash string `form:"block_hash"`
	Account   string `form:"account"`
	From      string `form:"from"`
	To        string `form:"to"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
	Limit     uint   `form:"limit"`

	startTime *time.Time
	endTime   *time.Time
}

func (s *TransactionSearch) Validate() error {
	if s.Type != "" {
		s.Type = strings.ToLower(s.Type)
		found := false
		for _, t := range model.TxTypes {
			if s.Type == t {
				found = true
				break
			}
		}
		if !found {
			return errors.New("invalid transaction type: " + s.Type)
		}
	}

	if t, err := parseTimeFilter(s.StartTime); err == nil {
		s.startTime = t
	} else {
		return errors.New("start time is invalid")
	}
	if t, err := parseTimeFilter(s.EndTime); err == nil {
		s.endTime = t
	} else {
		return errors.New("end time is invalid")
	}

	if s.BeforeID > 0 && s.AfterID > 0 {
		return errors.New("can't use both before/after ids")
	}

	if s.BlockHash != "" {
		if s.BeforeID > 0 || s.AfterID > 0 {
			return errors.New("can't use before/after with block hash")
		}
		if s.Height > 0 {
			return errors.New("can't use height with block hash")
		}
	}

	if s.startTime != nil && s.endTime != nil && s.endTime.Before(*s.startTime) {
		return errors.New("end time must be greater than start time")
	}

	if (s.startTime != nil || s.endTime != nil) && (s.BeforeID > 0 || s.AfterID > 0) {
		return errors.New("can't use time and ID filters together")
	}

	if s.Limit == 0 {
		s.Limit = 25
	}
	if s.Limit > 100 {
		s.Limit = 100
	}

	return nil
}

func parseTimeFilter(input string) (*time.Time, error) {
	if input == "" {
		return nil, nil
	}

	var t time.Time
	var err error

	if reDate.MatchString(input) {
		t, err = time.Parse("2006-01-02", input)
	} else {
		t, err = time.Parse(time.RFC3339, input)
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}
