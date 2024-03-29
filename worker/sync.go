package worker

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/figment-networks/mina-indexer/client/archive"
	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/client/staketab"
	"github.com/figment-networks/mina-indexer/config"
	"github.com/figment-networks/mina-indexer/indexing"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/mapper"
	"github.com/figment-networks/mina-indexer/model/types"
	"github.com/figment-networks/mina-indexer/store"
)

const unsafeBlockThreshold = 15

type SyncWorker struct {
	cfg            *config.Config
	db             *store.Store
	graphClient    *graph.Client
	archiveClient  *archive.Client
	staketabClient *staketab.Client
}

func NewSyncWorker(
	cfg *config.Config,
	db *store.Store,
	graphClient *graph.Client,
	archiveClient *archive.Client,
	staketabClient *staketab.Client,
) SyncWorker {
	return SyncWorker{
		cfg:            cfg,
		db:             db,
		graphClient:    graphClient,
		archiveClient:  archiveClient,
		staketabClient: staketabClient,
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
		Limit: 100,
	}
	if lastBlock != nil {
		blocksRequest.StartHeight = uint(lastBlock.Height + 1)
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
	if len(blocks) == 0 || (len(blocks) == 1 && lastBlock != nil && blocks[0].StateHash == lastBlock.Hash) {
		log.Info("no more blocks to process")
		return 0, nil
	}

	for _, block := range blocks {
		if err := w.processBlock(block.StateHash); err != nil {
			return 0, err
		}
	}

	lastBlock, err = w.db.Blocks.LastBlock()
	if err != nil {
		return 0, err
	}

	t := true
	blocksRequest = &archive.BlocksRequest{Canonical: &t}
	var limit uint = w.cfg.HistoricalLimit
	if (int(lastBlock.Height) - int(limit)) > 0 {
		blocksRequest.StartHeight = uint(lastBlock.Height+1) - limit
		blocksRequest.Limit = limit
	} else {
		blocksRequest.StartHeight = 0
		blocksRequest.Limit = uint(lastBlock.Height)
	}
	log.WithField("from", blocksRequest.StartHeight).Info("correcting canonical blocks")
	canonicalBlocks, err := w.archiveClient.Blocks(blocksRequest)
	if err != nil {
		return 0, err
	}
	for _, block := range canonicalBlocks {
		_, err := w.db.Blocks.FindByHash(block.StateHash)
		if err != nil {
			if err != store.ErrNotFound {
				return 0, err
			}
			if err := w.processBlock(block.StateHash); err != nil {
				return 0, err
			}
		}

		if err := w.db.Blocks.MarkBlocksOrphan(block.Height); err != nil {
			return 0, err
		}
		if err := w.db.Blocks.MarkBlockCanonical(block.StateHash); err != nil {
			return 0, err
		}
		if err := w.db.Transactions.MarkTransactionsOrphan(block.Height); err != nil {
			return 0, err
		}
		if err := w.db.Transactions.MarkTransactionsCanonical(block.StateHash); err != nil {
			return 0, err
		}
	}

	var unsafeBlocksStarting uint64
	if (int(lastBlock.Height) - int(limit)) > 0 {
		unsafeBlocksStarting = lastBlock.Height - unsafeBlockThreshold
	}
	log.WithField("from", unsafeBlocksStarting).Info("correcting canonical blocks and validators statistics")
	unsafeBlocks, err := w.db.Blocks.FindUnsafeBlocks(unsafeBlocksStarting)
	if err != nil {
		return 0, err
	}

	validatorKeys := map[string]string{}
	blockKeys := map[uint64]model.Block{}
	for _, b := range unsafeBlocks {
		_, ok := validatorKeys[b.Creator]
		if !ok {
			validatorKeys[b.Creator] = b.Creator
		}

		if !b.Canonical {
			continue
		}

		_, ok = blockKeys[b.Height]
		if !ok {
			blockKeys[b.Height] = b
		}
	}

	for _, block := range blockKeys {
		ts := block.Time
		buckets := []string{store.BucketHour, store.BucketDay}

		for _, bucket := range buckets {
			log.WithField("bucket", bucket).Debug("correcting chain stats")
			if err := w.db.Stats.CreateChainStats(bucket, ts); err != nil {
				return 0, err
			}

			log.WithField("bucket", bucket).Debug("creating validator stats")
			for _, key := range validatorKeys {
				if err := w.db.Stats.CreateValidatorStats(key, bucket, ts); err != nil {
					return 0, err
				}
			}
		}
	}

	safeCanonicalBlocksStarting := unsafeBlocksStarting - uint64(limit)
	lastCalculatedBlockReward, err := w.db.Blocks.FindLastCalculatedBlockReward(safeCanonicalBlocksStarting, unsafeBlocksStarting)
	if err != nil && err != store.ErrNotFound {
		return 0, err
	}
	if err == nil {
		safeCanonicalBlocksStarting = lastCalculatedBlockReward.Height
	}

	log.WithField("from", safeCanonicalBlocksStarting).WithField("to", unsafeBlocksStarting).Info("calculating rewards for safe canonical blocks")
	blocksForRewards, err := w.db.Blocks.FindNonCalculatedBlockRewards(safeCanonicalBlocksStarting, unsafeBlocksStarting)
	if err != nil {
		return 0, err
	}
	for _, block := range blocksForRewards {
		log.WithField("height", block.Height).WithField("hash", block.Hash).Info("calculating rewards")
		if err := indexing.RewardCalculation(w.db, block); err != nil {
			return 0, err
		}
	}

	// Dumping staging ledger is very expensive on the node and always times out.
	// It also locks the GraphQL endpoint, so no queries could be made during the dump.
	// We want to make this configurable rather then removing it completely.
	if !w.cfg.StagingLedgerDisabled {
		log.Info("processing staging ledger")
		if err := w.processStagingLedger(); err != nil {
			log.WithError(err).Error("staging ledger processing failed")
			// do not abort here
		}
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

	validatorEpochs := []model.ValidatorEpoch{}
	if graphBlock != nil {
		validatorEpochs, err = w.db.ValidatorsEpochs.GetValidatorEpochs(graphBlock.ProtocolState.ConsensusState.Epoch, "")
		if err != nil && err != store.ErrNotFound {
			return err
		}

		if len(validatorEpochs) == 0 && w.cfg.IdentityFile != "" {
			epoch, err := strconv.Atoi(graphBlock.ProtocolState.ConsensusState.Epoch)
			if err != nil {
				return err
			}

			uniqueValidatorEpochs := map[string]model.ValidatorEpoch{}

			err = indexing.ReadIdentityFile(w.cfg.IdentityFile, func(identity indexing.Identity) error {
				if identity.PublicKey == "" || identity.Fee == nil {
					return nil
				}

				validatorEpoch := model.ValidatorEpoch{
					AccountAddress: identity.PublicKey,
					ValidatorFee:   types.NewFloat64Float(*identity.Fee),
					Epoch:          epoch,
				}

				// check for duplicates, but only fail if the fee is different (ignore different Names since it isn't used for validatorEpochs)
				if prev, found := uniqueValidatorEpochs[identity.PublicKey]; found {
					if prev.ValidatorFee.String() != validatorEpoch.ValidatorFee.String() {
						return fmt.Errorf("duplicate validator PublicKey: %s with different fees %s vs %s", identity.PublicKey, prev.ValidatorFee.String(), validatorEpoch.ValidatorFee.String())
					} else {
						// if the fees are the same, don't add this epoch to the validatorEpochs slice
						return nil
					}
				}

				uniqueValidatorEpochs[validatorEpoch.AccountAddress] = validatorEpoch

				validatorEpochs = append(validatorEpochs, validatorEpoch)

				return nil
			})

			if err != nil {
				return err
			}

		}
	}

	log.
		WithField("hash", archiveBlock.StateHash).
		WithField("height", archiveBlock.Height).
		Debug("processing block")

	data, err := indexing.Prepare(archiveBlock, graphBlock, validatorEpochs)
	if err != nil {
		return err
	}

	// Fetch account details for all accounts seen in the block
	if err := w.fetchSeenAccounts(graphBlock, data); err != nil {
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

func (w SyncWorker) fetchSeenAccounts(graphBlock *graph.Block, data *indexing.Data) error {
	if graphBlock == nil {
		return nil
	}

	for _, accID := range data.AccountIDs() {
		acc, err := w.graphClient.GetAccount(accID)
		if err != nil {
			return err
		}

		// an account may be in the block, but doesn't actually exist because of a failed account creation transaction
		// e.g. Amount_insufficient_to_create_account
		// skip these since the account does not exist
		if acc == nil {
			log.WithField("account_id", accID).Info("skipping nonexistent account")
			continue
		}

		mappedAcc, err := mapper.Account(graphBlock, acc)
		if err != nil {
			return err
		}

		data.Accounts = append(data.Accounts, *mappedAcc)
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
		records, err := w.db.Staking.LedgerRecords(currentLedger.ID)
		if err != nil && err != store.ErrNotFound || len(records) > 0 {
			return nil, nil
		}
	}

	ledger, err := w.archiveClient.StakingLedger(archive.LedgerTypeCurrent)
	if err != nil {
		return nil, err
	}

	ledgerData, err := mapper.Ledger(tip, ledger)
	if err != nil {
		return nil, err
	}

	if currentLedger == nil {
		err = w.db.Staking.CreateLedger(ledgerData.Ledger)
		if err != nil {
			return nil, err
		}
	} else {
		ledgerData.Ledger = currentLedger
	}

	ledgerData.UpdateLedgerID()

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
