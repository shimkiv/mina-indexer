package worker

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/figment-networks/mina-indexer/coda"
	"github.com/figment-networks/mina-indexer/config"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/store"
)

func RunInit(cfg *config.Config, db *store.Store) error {
	n, err := db.Accounts.Count()
	if err != nil {
		return err
	}
	if n > 0 {
		return errors.New("accounts table should be empty before the genesis import")
	}

	genesis, err := coda.ReadGenesisFile(cfg.GenesisFile)
	if err != nil {
		return err
	}
	log.Info("genesis accounts found:", len(genesis.Accounts))

	// These are default start/end attributes and might not be correct
	// TODO: make the end height/time optional?
	height := uint64(0)
	now := time.Now()

	for _, a := range genesis.Accounts {
		acc := model.Account{
			PublicKey:   a.PK,
			Balance:     a.Balance,
			Delegate:    a.Delegate,
			StartHeight: height,
			StartTime:   now,
			LastHeight:  height,
			LastTime:    now,
		}

		log.WithField("pk", a.PK).Info("importing account")
		if err := db.Accounts.Create(&acc); err != nil {
			return err
		}
	}

	return nil
}
