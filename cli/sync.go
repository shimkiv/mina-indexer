package cli

import (
	"github.com/figment-networks/mina-indexer/client/archive"
	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/config"
	"github.com/figment-networks/mina-indexer/worker"
)

func runSync(cfg *config.Config) error {
	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	syncWorker := worker.NewSyncWorker(
		cfg, db,
		graph.NewDefaultClient(cfg.CodaEndpoint),
		archive.NewDefaultClient(cfg.ArchiveEndpoint),
	)

	_, err = syncWorker.Run()
	return err
}
