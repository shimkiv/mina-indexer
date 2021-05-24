package worker

import (
	"context"
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/figment-networks/mina-indexer/client/archive"
	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/config"
	"github.com/figment-networks/mina-indexer/indexing"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/mapper"
	"github.com/figment-networks/mina-indexer/store"
)

type SyncWorker struct {
	cfg           *config.Config
	db            *store.Store
	graphClient   *graph.Client
	archiveClient *archive.Client
}

func NewSyncWorker(
	cfg *config.Config,
	db *store.Store,
	graphClient *graph.Client,
	archiveClient *archive.Client,
) SyncWorker {
	return SyncWorker{
		cfg:           cfg,
		db:            db,
		graphClient:   graphClient,
		archiveClient: archiveClient,
	}
}

func (w SyncWorker) Run() (int, error) {
	log.Info("starting sync")

	status, err := w.checkNodeStatus()
	if err != nil {
		return 0, err
	}

	log.Info("processing staking ledger")
	_, err = w.processStakingLedger()
	if err != nil {
		return 0, err
	}

	log.Info("fetching the most recent indexed block")
	lastBlock, err := w.db.Blocks.Recent()
	if err != nil {
		log.Debug("latest indexed block is not found")
		if err != store.ErrNotFound {
			return 0, err
		}
		lastBlock = nil
	}

	blocksRequest := &archive.BlocksRequest{
		Limit:     100,
		Canonical: true,
	}
	if lastBlock != nil {
		blocksRequest.StartHeight = uint(lastBlock.Height)
	}

	log.
		WithField("start_height", blocksRequest.StartHeight).
		WithField("limit", blocksRequest.Limit).
		Info("fetching blocks from archive")

	blocks, err := w.archiveClient.Blocks(blocksRequest)
	if err != nil {
		return 0, err
	}

	// Remove the last block from the response. Last block in the chain is not 100%
	// canonical until we get the next block with the right parent hash.
	if lastBlock != nil && len(blocks) > 0 {
		blocks = blocks[0 : len(blocks)-1]
	}

	// Check if the only block we received is the most recent indexed one
	if len(blocks) == 0 || (len(blocks) == 1 && lastBlock != nil && blocks[0].StateHash == lastBlock.Hash) {
		log.Info("no more blocks to process")
		return 0, nil
	}

	for _, block := range blocks {
		if err := w.processBlock(block.StateHash); err != nil {
			return 0, err
		}
	}

	log.Info("processing staging ledger")
	if err := w.processStagingLedger(); err != nil {
		log.WithError(err).Error("staging ledger processing failed")
		// do not abort here
	}

	var lag int

	if len(blocks) > 0 {
		lag = status.HighestBlockLengthReceived - int(blocks[len(blocks)-1].Height)
		log.WithField("lag", lag).Info("sync finished")
	}

	return lag, err
}

func (w SyncWorker) processBlock(hash string) error {
	archiveBlock, err := w.archiveClient.Block(hash)
	if err != nil {
		return err
	}

	graphBlock, err := w.graphClient.GetBlock(hash)
	if err != nil {
		if !strings.Contains(err.Error(), "not found in transition frontier") {
			return err
		}

		log.WithError(err).Debug("graph block error")
		graphBlock = nil
	}

	log.
		WithField("hash", archiveBlock.StateHash).
		WithField("height", archiveBlock.Height).
		Debug("processing block")

	data, err := indexing.Prepare(archiveBlock, graphBlock)
	if err != nil {
		return err
	}

	if err := indexing.Import(w.db, data); err != nil {
		return err
	}

	if err := indexing.Finalize(w.db, data); err != nil {
		return err
	}

	return nil
}

func (w SyncWorker) checkNodeStatus() (*graph.DaemonStatus, error) {
	log.Debug("fetching node status")
	status, err := w.graphClient.GetDaemonStatus(context.Background())
	if err != nil {
		return nil, err
	}

	log.
		WithField("status", status.SyncStatus).
		Debug("current node status")

	switch status.SyncStatus {
	case graph.SyncStatusOffline:
		return nil, errors.New("node is offline")
	case graph.SyncStatusConnecting:
		return nil, errors.New("node is connecting")
	case graph.SyncStatusBootstrap:
		return nil, errors.New("node is bootstrapping")
	}

	return status, nil
}

func (w SyncWorker) processStakingLedger() (*mapper.LedgerData, error) {
	tip, err := w.graphClient.ConsensusTip()
	if err != nil {
		return nil, err
	}

	var epoch int
	fmt.Sscanf(tip.ProtocolState.ConsensusState.Epoch, "%d", &epoch)

	// Find ledger for current epoch. Ledger only changes once per epoch.
	currentLedger, err := w.db.Staking.FindLedger(epoch)
	if err != nil && err != store.ErrNotFound {
		return nil, err
	}

	// We already have current epoch ledger, no need to import it.
	if currentLedger != nil && currentLedger.EntriesCount > 0 {
		return nil, nil
	}

	ledger, err := w.archiveClient.StakingLedger(archive.LedgerTypeCurrent)
	if err != nil {
		return nil, err
	}

	ledgerData, err := mapper.Ledger(tip, ledger)
	if err != nil {
		return nil, err
	}

	err = w.db.Staking.CreateLedger(ledgerData.Ledger)
	if err != nil {
		return nil, err
	}
	ledgerData.UpdateLedgerID()

	err = ledgerData.SetWeights()
	if err != nil {
		return nil, err
	}

	err = w.db.Staking.CreateLedgerEntries(ledgerData.Entries)
	if err != nil {
		return nil, err
	}

	return ledgerData, nil
}

func (w SyncWorker) processStagingLedger() error {
	tip, err := w.graphClient.ConsensusTip()
	if err != nil {
		log.WithError(err).Error("consensus tip fetch failed")
		return err
	}

	block, err := w.graphClient.GetBlock(tip.StateHash)
	if err != nil {
		log.WithError(err).Error("block fetch failed")
		return err
	}

	ledger, err := w.archiveClient.StakingLedger(archive.LedgerTypeStaged)
	if err != nil {
		log.WithError(err).Error("staged ledger fetch failed")
		return err
	}

	accounts := []model.Account{}
	for _, entry := range ledger {
		account, err := mapper.AccountFromStagedLedger(block, &entry)
		if err != nil {
			log.WithError(err).Error("account init failed")
			continue
		}
		accounts = append(accounts, *account)
	}

	return w.db.Accounts.Import(accounts)
}
