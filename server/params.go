package server

import "time"

type rewardsParams struct {
	From     time.Time `form:"from" binding:"required" time_format:"2006-01-02"`
	To       time.Time `form:"to" binding:"required" time_format:"2006-01-02"`
	Interval string    `form:"interval" binding:"required" `
}

type delegatorRewardsParams struct {
	rewardsParams
	ValidatorId string `form:"validator_id" binding:"-" `
}
