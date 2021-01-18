package cli

import (
	"errors"

	"github.com/figment-networks/mina-indexer/config"
	"github.com/figment-networks/mina-indexer/worker"
)

func startInit(cfg *config.Config) error {
	if cfg.GenesisFile == "" {
		return errors.New("genesis file is not provided")
	}

	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	return worker.RunInit(cfg, db)
}
