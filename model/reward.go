package model

import (
	"time"

	"github.com/figment-networks/mina-indexer/model/types"
)

// BlockReward contains the reward details earned at a specific height
type BlockReward struct {
	ID           string       `json:"-"`
	OwnerAccount string       `json:"owner_account"`
	Delegate     *string      `json:"delegate"`
	Epoch        int          `json:"epoch"`
	TimeBucket   time.Time    `json:"time_bucket"`
	Reward       types.Amount `json:"reward"`
	OwnerType    string       `json:"owner_type"`
}

// String returns account text representation
func (br BlockReward) String() string {
	return br.OwnerAccount
}

type RewardsSummary struct {
	Interval string       `json:"interval"`
	Epoch    string       `json:"epoch,omitempty"`
	Delegate string       `json:"delegate,omitempty"`
	Amount   types.Amount `json:"amount"`
}

type TimeInterval uint
type RewardOwnerType string

const (
	TimeIntervalDaily TimeInterval = iota
	TimeIntervalMonthly
	TimeIntervalYearly

	RewardOwnerTypeValidator RewardOwnerType = "validator"
	RewardOwnerTypeDelegator RewardOwnerType = "delegator"
)

var (
	TimeIntervalTypes = map[string]TimeInterval{
		"daily":   TimeIntervalDaily,
		"monthly": TimeIntervalMonthly,
		"yearly":  TimeIntervalYearly,
	}

	RewardOwnerTypes = map[string]RewardOwnerType{
		"validator": RewardOwnerTypeValidator,
		"delegator": RewardOwnerTypeDelegator,
	}
)

func GetTypeForTimeInterval(s string) (TimeInterval, bool) {
	t, ok := TimeIntervalTypes[s]
	return t, ok
}

func (k TimeInterval) String() string {
	switch k {
	case TimeIntervalDaily:
		return "YYYY-MM-DD"
	case TimeIntervalMonthly:
		return "YYYY-MM"
	case TimeIntervalYearly:
		return "YYYY"
	default:
		return "unknown"
	}
}

func GetTypeForRewardOwnerType(s string) (RewardOwnerType, bool) {
	t, ok := RewardOwnerTypes[s]
	return t, ok
}
