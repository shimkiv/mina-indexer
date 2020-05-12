package cli

import (
	"net/http"

	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/config"
	"github.com/figment-networks/coda-indexer/pipeline"
)

func runSync(cfg *config.Config) error {
	client := coda.NewClient(http.DefaultClient, cfg.CodaEndpoint)

	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	return pipeline.NewSync(cfg, db, client).Execute()
}
