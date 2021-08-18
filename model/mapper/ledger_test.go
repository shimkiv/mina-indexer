package mapper

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/figment-networks/mina-indexer/client/archive"
	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/model/types"
	"github.com/stretchr/testify/assert"
)

func readLedgerFromFile(path string) []archive.StakingInfo {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	result := []archive.StakingInfo{}
	if err := json.NewDecoder(f).Decode(&result); err != nil {
		panic(err)
	}

	return result
}

func TestLedger(t *testing.T) {
	entries := readLedgerFromFile("../../test/fixtures/ledger.json")

	block := &graph.Block{
		ProtocolState: &graph.ProtocolState{
			ConsensusState: &graph.ConsensusState{
				Epoch: "10",
			},
		},
	}

	ledgerData, err := Ledger(block, entries)
	ledger := ledgerData.Ledger

	assert.NoError(t, err)
	assert.Equal(t, 3, ledger.EntriesCount)
	assert.Equal(t, 10, ledger.Epoch)
	assert.Equal(t, 3, len(ledgerData.Entries))

	entry := ledgerData.Entries[0]
	assert.Equal(t, types.NewFloatAmount("4651").String(), entry.TimingInitialMinimumBalance.String())
	assert.Equal(t, 86400, *entry.TimingCliffTime)
	assert.Equal(t, types.NewFloatAmount("4651").String(), entry.TimingCliffAmount.String())
	assert.Equal(t, 1, *entry.TimingVestingPeriod)
	assert.Equal(t, types.NewFloatAmount("0").String(), entry.TimingVestingIncrement.String())

	entry = ledgerData.Entries[1]
	assert.Equal(t, types.NewFloatAmount("66000").String(), entry.TimingInitialMinimumBalance.String())
	assert.Equal(t, 172800, *entry.TimingCliffTime)
	assert.Equal(t, types.NewFloatAmount("16500").String(), entry.TimingCliffAmount.String())
	assert.Equal(t, 1, *entry.TimingVestingPeriod)
	assert.Equal(t, types.NewFloatAmount("0.095486111").String(), entry.TimingVestingIncrement.String())

	entry = ledgerData.Entries[2]
	assert.Equal(t, types.NewFloatAmount("15277795.5092651").String(), entry.TimingInitialMinimumBalance.String())
	assert.Equal(t, 43200, *entry.TimingCliffTime)
	assert.Equal(t, types.NewFloatAmount("0").String(), entry.TimingCliffAmount.String())
	assert.Equal(t, 1, *entry.TimingVestingPeriod)
	assert.Equal(t, types.NewFloatAmount("70.730534765").String(), entry.TimingVestingIncrement.String())
}
