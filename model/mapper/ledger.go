package mapper

import (
	"fmt"
	"time"

	"github.com/figment-networks/mina-indexer/client/archive"
	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/types"
	"github.com/figment-networks/mina-indexer/model/util"
)

type LedgerData struct {
	Ledger  *model.Ledger
	Entries []model.LedgerEntry
}

func (data *LedgerData) UpdateLedgerID() {
	for idx := range data.Entries {
		data.Entries[idx].LedgerID = data.Ledger.ID
	}
}

func (data *LedgerData) SetWeights() error {
	for idx := range data.Entries {
		w, err := util.CalculateWeight(data.Entries[idx].Balance, data.Ledger.StakedAmount)
		if err != nil {
			return err
		}
		data.Entries[idx].Weight = w
	}
	return nil
}

func Ledger(tip *graph.Block, records []archive.StakingInfo) (*LedgerData, error) {
	ledgerRecord := &model.Ledger{
		EntriesCount:      len(records),
		Time:              time.Now(),
		DelegationsAmount: types.NewInt64Amount(0),
		StakedAmount:      types.NewInt64Amount(0),
	}
	fmt.Sscanf(tip.ProtocolState.ConsensusState.Epoch, "%d", &ledgerRecord.Epoch)

	entries := []model.LedgerEntry{}

	for _, record := range records {
		balance := types.NewFloatAmount(record.Balance)

		entry := model.LedgerEntry{
			LedgerID:                    ledgerRecord.ID,
			PublicKey:                   record.Pk,
			Delegate:                    record.Delegate,
			Balance:                     balance,
			Delegation:                  record.Pk != record.Delegate,
			TimingInitialMinimumBalance: types.Amount{},
			TimingCliffAmount:           types.Amount{},
		}

		ledgerRecord.StakedAmount = ledgerRecord.StakedAmount.Add(balance)

		if entry.Delegation {
			ledgerRecord.DelegationsCount++
			ledgerRecord.DelegationsAmount = ledgerRecord.DelegationsAmount.Add(balance)
		}

		if timing := record.Timing; timing != nil {
			cliffTime, _ := util.ParseInt(timing.CliffTime)
			vestingIncrement, _ := util.ParseInt(timing.VestingIncrement)
			vestingPeriod, _ := util.ParseInt(timing.VestingPeriod)

			entry.TimingInitialMinimumBalance = types.NewFloatAmount(timing.InitialMinimumBalance)
			entry.TimingCliffAmount = types.NewFloatAmount(timing.CliffAmount)
			entry.TimingCliffTime = &cliffTime
			entry.TimingVestingIncrement = &vestingIncrement
			entry.TimingVestingPeriod = &vestingPeriod
		}

		entries = append(entries, entry)
	}

	return &LedgerData{
		Ledger:  ledgerRecord,
		Entries: entries,
	}, nil
}
