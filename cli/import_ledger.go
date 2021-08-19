package cli

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"github.com/figment-networks/mina-indexer/client/archive"
	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/config"
	"github.com/figment-networks/mina-indexer/model/mapper"
)

func importLedger(cfg *config.Config) error {
	store, err := initStore(cfg)
	if err != nil {
		return err
	}

	ledgerFile := os.Getenv("LEDGER_FILE")
	if ledgerFile == "" {
		return errors.New("LEDGER_FILE env var is not provided")
	}

	epochVal := os.Getenv("LEDGER_EPOCH")
	if epochVal == "" {
		return errors.New("LEDGER_EPOCH env var is not provided")
	}

	epoch, err := strconv.Atoi(epochVal)
	if err != nil {
		return err
	}

	rawData, err := decodeLedgerData(ledgerFile)
	if err != nil {
		return err
	}

	tip := &graph.Block{
		ProtocolState: &graph.ProtocolState{
			ConsensusState: &graph.ConsensusState{
				Epoch: epochVal,
			},
		},
	}

	ledger, err := mapper.Ledger(tip, rawData)
	if err != nil {
		return nil
	}

	if err := store.Staking.DeleteEpochLedger(epoch); err != nil {
		return err
	}

	if err := store.Staking.CreateLedger(ledger.Ledger); err != nil {
		return err
	}
	ledger.UpdateLedgerID()

	return store.Staking.CreateLedgerEntries(ledger.Entries)
}

func decodeLedgerData(path string) ([]archive.StakingInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := []archive.StakingInfo{}
	return result, json.NewDecoder(f).Decode(&result)
}
