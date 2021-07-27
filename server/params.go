package server

import (
	"errors"
	"time"

	"github.com/figment-networks/mina-indexer/model"
)

type rewardsParams struct {
	From            time.Time `form:"from" time_format:"2006-01-02"`
	To              time.Time `form:"to" time_format:"2006-01-02"`
	ValidatorId     string    `form:"validator_id"`
	RewardOwnerType string    `form:"owner_type" binding:"required" `
	Interval        string    `form:"interval" binding:"required" `
}

func (p *rewardsParams) Validate() error {
	var ok bool
	if _, ok = model.GetTypeForTimeInterval(p.Interval); !ok {
		return errors.New("time interval type is wrong")
	}

	var ownerType model.RewardOwnerType
	if ownerType, ok = model.GetTypeForRewardOwnerType(p.RewardOwnerType); !ok {
		return errors.New("owner type is wrong")
	}

	if ownerType == model.RewardOwnerTypeDelegator && p.ValidatorId == "" {
		return errors.New("validator id should be defined for delegator reward")
	}

	return nil
}
