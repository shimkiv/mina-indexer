package worker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/config"
	"github.com/figment-networks/coda-indexer/indexing"
	"github.com/figment-networks/coda-indexer/model/mapper"
	"github.com/figment-networks/coda-indexer/store"
)

func RunSync(cfg *config.Config, db *store.Store, client *coda.Client) (int, error) {
	log.Debug("checking daemon status")
	status, err := checkNodeStatus(client)
	if err != nil {
		return 0, err
	}

	var blocks []coda.Block
	var lag int

	log.Info("fetching blocks")
	defer func() {
		log.WithField("lag", lag).Info("done fetching")
	}()

	block, err := db.Blocks.Recent()
	if err != nil {
		if err == store.ErrNotFound {
			blocks, err = client.GetFirstBlocks(10)
		} else {
			return -1, err
		}
	}
	lag = *status.BlockchainLength - int(block.Height)

	if len(blocks) == 0 {
		blocks, err = client.GetNextBlocks(block.Hash, 10)
		if err != nil {
			return lag, err
		}
	}

	n, err := processBlocks(cfg, db, status, blocks)
	lag -= n

	return lag, err
}

func checkNodeStatus(client *coda.Client) (*coda.DaemonStatus, error) {
	log.Debug("fetching node status")
	status, err := client.GetDaemonStatus()
	if err != nil {
		return nil, err
	}
	log.
		WithField("status", status.SyncStatus).
		Debug("current node status")

	switch status.SyncStatus {
	case coda.SyncStatusOffline:
		return nil, errors.New("node is offline")
	case coda.SyncStatusConnecting:
		return nil, errors.New("node is connecting")
	case coda.SyncStatusBootstrap:
		return nil, errors.New("node is bootstrapping")
	}

	return status, nil
}

func processBlocks(cfg *config.Config, db *store.Store, status *coda.DaemonStatus, blocks []coda.Block) (int, error) {
	imported := 0
	total := len(blocks)

	defer func() {
		log.
			WithFields(log.Fields{"fetched": total, "imported": imported}).
			Info("done processing")
	}()

	for _, block := range blocks {
		if cfg.DumpDir != "" {
			dumpBlock(&block, cfg.DumpDir)
		}

		done, err := processBlock(db, status, &block)
		if err != nil {
			return imported, err
		}
		if done {
			imported++
		}
	}

	return imported, nil
}

func processBlock(db *store.Store, status *coda.DaemonStatus, block *coda.Block) (bool, error) {
	_, err := db.Blocks.FindByHash(block.StateHash)
	if err == nil {
		log.WithField("hash", block.StateHash).Debug("skipping already existing block")
		return false, nil
	}
	if err != store.ErrNotFound {
		return false, err
	}

	log.WithField("hash", block.StateHash).Info("processing block")
	data, err := indexing.Prepare(status, block)
	if err != nil {
		return false, err
	}

	if err := indexing.Import(db, data); err != nil {
		return false, err
	}

	if err := indexing.Finalize(db, data); err != nil {
		return false, err
	}

	return true, nil
}

func dumpBlock(block *coda.Block, dir string) error {
	time := mapper.BlockTime(block)

	savePath := fmt.Sprintf("%v/%v/%v_%v.json",
		dir,
		time.Format("2006-01-02"),
		block.ProtocolState.ConsensusState.BlockHeight,
		block.StateHash[0:64],
	)

	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		log.WithError(err).Error("dump dir creation failed")
		panic(err)
	}

	data, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		log.WithError(err).Error("block dump json failed")
		return err
	}

	log.WithField("path", savePath).Debug("saving block data")
	if err := ioutil.WriteFile(savePath, data, 0755); err != nil {
		log.WithError(err).Error("block write failed")
		return err
	}

	return nil
}
