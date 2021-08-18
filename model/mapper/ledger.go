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
			cliffTime, err := util.ParseInt(timing.CliffTime)
			if err != nil {
				return nil, err
			}

			vestingPeriod, err := util.ParseInt(timing.VestingPeriod)
			if err != nil {
				return nil, err
			}

			entry.TimingInitialMinimumBalance = types.NewFloatAmount(timing.InitialMinimumBalance)
			entry.TimingCliffAmount = types.NewFloatAmount(timing.CliffAmount)
			entry.TimingCliffTime = &cliffTime
			entry.TimingVestingIncrement = types.NewFloatAmount(timing.VestingIncrement)
			entry.TimingVestingPeriod = &vestingPeriod
		}

		entries = append(entries, entry)
	}

	return &LedgerData{
		Ledger:  ledgerRecord,
		Entries: entries,
	}, nil
}
