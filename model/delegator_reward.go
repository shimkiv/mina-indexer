package model

import (
	"time"

	"github.com/figment-networks/mina-indexer/model/types"
)

// DelegatorBlockReward contains the delegator reward details earned at a specific height
type DelegatorBlockReward struct {
	ID          string       `json:"-"`
	PublicKey   string       `json:"public_key"`
	Delegate    *string      `json:"delegate"`
	BlockHeight uint64       `json:"block_height"`
	BlockTime   time.Time    `json:"block_time"`
	Reward      types.Amount `json:"reward"`
}

// String returns account text representation
func (dbr DelegatorBlockReward) String() string {
	return dbr.PublicKey
}

type RewardsSummary struct {
	Interval string       `json:"interval"`
	Amount   types.Amount `json:"amount"`
}

type TimeInterval uint

const (
	TimeIntervalDaily TimeInterval = iota
	TimeIntervalMonthly
	TimeIntervalYearly
)

var (
	TimeIntervalTypes = map[string]TimeInterval{
		"daily":   TimeIntervalDaily,
		"monthly": TimeIntervalMonthly,
		"yearly":  TimeIntervalYearly,
	}
)

func GetTypeForTimeInterval(s string) (TimeInterval, bool) {
	t, ok := TimeIntervalTypes[s]
	return t, ok
}

func (k TimeInterval) String() string {
	switch k {
	case TimeIntervalDaily:
		return "DD-MM-YYYY"
	case TimeIntervalMonthly:
		return "MM-YYYY"
	case TimeIntervalYearly:
		return "YYYY"
	default:
		return "unknown"
	}
}
