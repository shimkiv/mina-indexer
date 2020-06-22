package cli

import (
	"log"

	"github.com/figment-networks/coda-indexer/config"
	"github.com/figment-networks/coda-indexer/server"
)

func startServer(cfg *config.Config) error {
	server.SetGinDefaults(cfg)

	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Println("Starting server on", cfg.ListenAddr())
	return server.New(db, cfg).Run(cfg.ListenAddr())
}
