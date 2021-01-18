package indexing

import (
	log "github.com/sirupsen/logrus"

	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/store"
)

// Import creates new database records for the chain data
func Import(db *store.Store, data *Data) error {
	log.WithField("count", 1).Debug("creating accounts")
	if err := db.Accounts.Import(data.Accounts); err != nil {
		return err
	}

	log.Debug("creating block")
	if err := db.Blocks.Create(data.Block); err != nil {
		return err
	}

	log.WithField("count", 1).Debug("creating validators")
	if err := db.Validators.Import([]model.Validator{*data.Validator}); err != nil {
		return err
	}

	log.WithField("count", len(data.Transactions)).Debug("creating transactions")
	for _, t := range data.Transactions {
		et, _ := db.Transactions.FindByHash(t.Hash)
		if et.ID > 0 {
			log.WithField("id", et.Hash).Debug("transaction already exists")
			continue
		}

		if err := db.Transactions.Create(&t); err != nil {
			return err
		}
	}

	log.WithField("count", len(data.FeeTransfers)).Debug("creating fee transfers")
	if err := db.FeeTransfers.Import(data.FeeTransfers); err != nil {
		return err
	}

	log.WithField("count", len(data.Snarkers)).Debug("creating snarkers")
	if err := db.Snarkers.Import(data.Snarkers); err != nil {
		return err
	}

	log.WithField("count", len(data.SnarkJobs)).Debug("creating snarkjobs")
	if err := db.Jobs.Import(data.SnarkJobs); err != nil {
		return err
	}

	return nil
}
