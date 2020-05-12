package pipeline

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/config"
	"github.com/figment-networks/coda-indexer/model"
	"github.com/figment-networks/coda-indexer/model/mapper"
	"github.com/figment-networks/coda-indexer/model/util"
	"github.com/figment-networks/coda-indexer/store"
)

// Sync pipeline runs the database sync
type Sync struct {
	cfg    *config.Config
	db     *store.Store
	client *coda.Client

	currentHeight int64
	currentHash   string

	errors    []error
	status    *coda.DaemonStatus
	report    *model.Report
	syncables []*model.Syncable
}

// NewSync returns a new sync pipeline
func NewSync(cfg *config.Config, db *store.Store, client *coda.Client) *Sync {
	return &Sync{
		cfg:       cfg,
		db:        db,
		client:    client,
		errors:    []error{},
		syncables: []*model.Syncable{},
	}
}

// Execute runs the pipeline steps
func (s *Sync) Execute() error {
	return runChain(
		func() error { return s.checkDaemonStatus() },
		func() error { return s.getCurrentHeight() },
		func() error { return s.createReport() },
		func() error { return s.startReport() },
		func() error { return s.createSyncables() },
		func() error { return s.processSyncables() },
		func() error { return s.finishReport() },
	)
}

func (s *Sync) checkDaemonStatus() error {
	status, err := s.client.GetDaemonStatus()
	if err != nil {
		return err
	}
	if status.SyncStatus != coda.SyncStatusSynced {
		return errors.New("coda daemon is not synced")
	}

	s.status = status
	return nil
}

func (s *Sync) getCurrentHeight() error {
	// Figure out the most recent blocks
	block, err := s.getRecentBlock()
	if err != nil {
		return err
	}

	// Parse out current block height
	height, err := util.ParseInt64(block.ProtocolState.ConsensusState.BlockHeight)
	if err != nil {
		return err
	}

	// Set the current block information
	// NOTE: coda's chain info from API starts with height = 2
	s.currentHeight = height - 1

	return nil
}

func (s *Sync) getRecentBlock() (*coda.Block, error) {
	var block *coda.Block

	// Get the latest processed block syncable
	recent, err := s.db.Syncables.GetMostRecent(model.SyncableTypeBlock)
	if err != nil {
		if err != store.ErrNotFound {
			return nil, err
		}

		// No syncables found, fetch the first block from the node itself
		b, err := s.client.GetFirstBlock()
		if err != nil {
			return nil, err
		}
		block = b
	} else {
		// Parse out last available raw block
		if err := recent.Decode(&block); err != nil {
			return nil, err
		}

		// We will be fetching the next block after this hash
		s.currentHash = block.StateHash
	}

	return block, nil
}

func (s *Sync) createReport() error {
	startHeight := s.currentHeight
	if s.currentHash != "" {
		startHeight++
	}

	s.report = &model.Report{
		StartHeight: startHeight,
		EndHeight:   startHeight + 1,
		State:       model.ReportStatePending,
	}

	return s.db.Reports.Create(s.report)
}

func (s *Sync) startReport() error {
	s.report.State = model.ReportStateRunning
	return s.db.Reports.Update(s.report)
}

func (s *Sync) finishReport() error {
	s.report.Complete(0, 0, nil, nil)
	return s.db.Reports.Update(s.report)
}

func (s *Sync) createSyncables() error {
	exists, err := s.db.Syncables.Exists(model.SyncableTypeBlock, s.report.StartHeight)
	if err != nil {
		return err
	}
	if exists {
		log.Println("block syncable at height", s.report.StartHeight, "already exists")
		return nil
	}

	block, err := s.client.GetNextBlock(s.currentHash)
	if err != nil {
		return err
	}
	if block == nil {
		log.Println("no new blocks found")
		return nil
	}

	data, err := json.Marshal(block)
	if err != nil {
		return err
	}

	syncable := &model.Syncable{
		ReportID: s.report.ID,
		Height:   s.report.StartHeight,
		Time:     time.Now(),
		Type:     model.SyncableTypeBlock,
		Data:     data,
	}

	log.Println("creating syncable:", syncable)
	if err := s.db.Syncables.Create(syncable); err != nil {
		return err
	}
	s.syncables = append(s.syncables, syncable)

	return nil
}

func (s *Sync) processSyncables() error {
	if len(s.syncables) == 0 {
		log.Println("no syncables to process")
		return nil
	}

	for _, syncable := range s.syncables {
		if err := s.processSyncable(syncable); err != nil {
			log.Println("processing syncable:", syncable)
			return err
		}
	}

	return nil
}

func (s *Sync) processSyncable(syncable *model.Syncable) error {
	codaBlock := coda.Block{}

	err := runChain(
		func() error { return syncable.Decode(&codaBlock) },
		func() error { return s.createBlock(codaBlock) },
		func() error { return s.createState(codaBlock) },
		func() error { return s.createAccounts(codaBlock) },
		func() error { return s.createTransactions(codaBlock) },
	)
	if err == nil {
		err = s.db.Syncables.MarkProcessed(syncable)
	}
	return err
}

func (s *Sync) createBlock(block coda.Block) error {
	b, err := mapper.Block(block)
	if err != nil {
		return err
	}
	b.AppVersion = s.status.CommitID

	log.Println("creating block:", b)
	return s.db.Blocks.CreateIfNotExist(b)
}

func (s *Sync) createState(block coda.Block) error {
	state, err := mapper.State(block)
	if err != nil {
		return err
	}

	log.Println("creating state:", state)
	return s.db.States.CreateIfNotExists(state)
}

func (s *Sync) createAccounts(block coda.Block) error {
	accounts, err := mapper.Accounts(block)
	if err != nil {
		return err
	}

	for _, acc := range accounts {
		log.Println("creating account:", acc)
		if err := s.db.Accounts.CreateOrUpdate(&acc); err != nil {
			return err
		}
	}

	return nil
}

func (s *Sync) createTransactions(block coda.Block) error {
	transactions, err := mapper.Transactions(block)
	if err != nil {
		return err
	}
	for _, t := range transactions {
		log.Println("creating transaction", t)
		if err := s.db.Transactions.CreateIfNotExist(&t); err != nil {
			return err
		}
	}
	return nil
}
