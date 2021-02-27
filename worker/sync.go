package worker

import (
	"context"
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/figment-networks/mina-indexer/client/archive"
	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/config"
	"github.com/figment-networks/mina-indexer/indexing"
	"github.com/figment-networks/mina-indexer/model/types"
	"github.com/figment-networks/mina-indexer/store"
)

type SyncWorker struct {
	cfg           *config.Config
	db            *store.Store
	graphClient   *graph.Client
	archiveClient *archive.Client
}

func NewSyncWorker(cfg *config.Config, db *store.Store, graphClient *graph.Client, archiveClient *archive.Client) SyncWorker {
	return SyncWorker{
		cfg:           cfg,
		db:            db,
		graphClient:   graphClient,
		archiveClient: archiveClient,
	}
}

func (w SyncWorker) Run() (int, error) {
	status, err := w.checkNodeStatus()
	if err != nil {
		return 0, err
	}

	log.Debug("fetching the most recent indexed block")
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

	// Check if the only block we received is the most recent indexed one
	if len(blocks) == 1 && lastBlock != nil && blocks[0].StateHash == lastBlock.Hash {
		log.Info("no more blocks to process")
		return 0, nil
	}

	for _, block := range blocks {
		if err := w.processBlock(block.StateHash); err != nil {
			return 0, err
		}
	}

	var lag int

	if len(blocks) > 0 {
		lag = status.HighestBlockLengthReceived - int(blocks[len(blocks)-1].Height)
		log.WithField("lag", lag).Info("sync finished")
	}

	w.processStakingLedger()

	return lag, nil
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

func (w SyncWorker) processStakingLedger() error {
	ledger, err := w.archiveClient.StakingLedger("current")
	if err != nil {
		return err
	}

	totalStake := map[string]types.Amount{}

	for _, item := range ledger {
		balance := types.NewFloatAmount(item.Balance)

		if _, ok := totalStake[item.Delegate]; !ok {
			totalStake[item.Delegate] = types.NewAmount("0")
		} else {
			totalStake[item.Delegate] = totalStake[item.Delegate].Add(balance)
		}
	}

	for pk, stake := range totalStake {
		if err := w.db.Validators.UpdateStake(pk, stake); err != nil {
			return err
		}
	}

	return nil
}
