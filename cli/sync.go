package cli

import (
	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/config"
	"github.com/figment-networks/coda-indexer/worker"
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

	return worker.RunSync(cfg, db, client)
}
