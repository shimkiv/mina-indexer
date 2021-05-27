package model

import (
	"time"

	"github.com/figment-networks/mina-indexer/model/types"
)

// BlockReward contains the reward details earned at a specific height
type BlockReward struct {
	ID              string       `json:"-"`
	PublicKey       string       `json:"public_key"`
	Delegate        string       `json:"delegate"`
	BlockHeight     uint64       `json:"block_height"`
	BlockTime       time.Time    `json:"block_time"`
	Reward          types.Amount `json:"reward"`
	RewardOwnerType string       `json:"reward_owner_type"`
}

// String returns account text representation
func (dbr BlockReward) String() string {
	return dbr.PublicKey
}

type RewardsSummary struct {
	Interval string       `json:"interval"`
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
		return "DD-MM-YYYY"
	case TimeIntervalMonthly:
		return "MM-YYYY"
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
