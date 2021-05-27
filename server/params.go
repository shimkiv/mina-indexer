package server

import "time"

type rewardsParams struct {
	From            time.Time `form:"from" binding:"required" time_format:"2006-01-02"`
	To              time.Time `form:"to" binding:"required" time_format:"2006-01-02"`
	ValidatorId     string    `form:"validator_id" binding:"-" `
	RewardOwnerType string    `form:"owner_type" binding:"required" `
	Interval        string    `form:"interval" binding:"required" `
}
