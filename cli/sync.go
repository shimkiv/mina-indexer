package cli

import (
	"github.com/figment-networks/mina-indexer/coda"
	"github.com/figment-networks/mina-indexer/config"
	"github.com/figment-networks/mina-indexer/worker"
)

func runSync(cfg *config.Config) error {
	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	client := coda.NewDefaultClient(cfg.CodaEndpoint)
	if cfg.LogLevel == "debug" {
		client.SetDebug(true)
	}

	_, err = worker.RunSync(cfg, db, client)
	return err
}
