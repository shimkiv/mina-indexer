package model

import (
	"github.com/figment-networks/mina-indexer/model/types"
	"time"
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
