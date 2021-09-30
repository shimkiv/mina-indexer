package server

import (
	"time"

	"github.com/figment-networks/mina-indexer/model"
)

type HealthResponse struct {
	Healthy bool `json:"healthy"`
}

type StatusResponse struct {
	AppName         string    `json:"app_name"`
	AppVersion      string    `json:"app_version"`
	GitCommit       string    `json:"git_commit"`
	GoVersion       string    `json:"go_version"`
	NodeVersion     string    `json:"node_version,omitempty"`
	NodeStatus      string    `json:"node_status,omitempty"`
	NodeError       bool      `json:"node_error"`
	SyncStatus      string    `json:"sync_status"`
	LastBlockTime   time.Time `json:"last_block_time"`
	LastBlockHeight uint64    `json:"last_block_height"`
}

type HeightResponse struct {
	Height uint64    `json:"height"`
	Time   time.Time `json:"time"`
}

type BlockResponse struct {
	Block        *model.Block        `json:"block"`
	Creator      *model.Account      `json:"creator"`
	Transactions []model.Transaction `json:"transactions"`
	SnarkJobs    []model.SnarkJob    `json:"snark_jobs"`
}

type ValidatorResponse struct {
	Validator   *model.Validator      `json:"validator"`
	Account     *model.Account        `json:"account"`
	Delegations []model.Delegation    `json:"delegations"`
	Stats       []model.ValidatorStat `json:"stats"`
	StatsHourly []model.ValidatorStat `json:"stats_hourly"`
	StatsDaily  []model.ValidatorStat `json:"stats_daily"`
}

type LedgerRequest struct {
	Epoch *int `form:"epoch"`
}

type LedgerResponse struct {
	Ledger  *model.Ledger       `json:"ledger"`
	Records []model.LedgerEntry `json:"entries"`
}

type SnarkerResponse struct {
	Snarker model.Snarker         `json:"snarker"`
	Stats   []model.SnarkerStat   `json:"stats"`
	Fees    []model.SnarkerJobFee `json:"job_fees"`
}
