package cli

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/config"
	"github.com/figment-networks/coda-indexer/pipeline"
	"github.com/figment-networks/coda-indexer/store"
)

func startSync(cfg *config.Config, db *store.Store) error {
	log.Println("sync will run every", cfg.SyncInterval)
	duration, err := time.ParseDuration(cfg.SyncInterval)
	if err != nil {
		return err
	}

	log.Println("using coda endpoint", cfg.CodaEndpoint)
	client := coda.NewClient(http.DefaultClient, cfg.CodaEndpoint)

	for range time.Tick(duration) {
		log.Println("starting sync")
		if err := pipeline.NewSync(cfg, db, client).Execute(); err != nil {
			log.Println("sync error:", err)
		}
	}

	return nil
}

func startCleanup(cfg *config.Config, db *store.Store) error {
	log.Println("cleanup will run every", cfg.CleanupInterval)
	duration, err := time.ParseDuration(cfg.CleanupInterval)
	if err != nil {
		return err
	}

	for range time.Tick(duration) {
		log.Println("starting cleanup")
		if err := pipeline.NewCleanup(cfg, db).Execute(); err != nil {
			log.Println("clenaup error:", err)
		}
	}
	return nil
}

func startWorker(cfg *config.Config) error {
	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		if err := startSync(cfg, db); err != nil {
			log.Println(err)
		}
		wg.Done()
	}()

	go func() {
		if err := startCleanup(cfg, db); err != nil {
			log.Println(err)
		}
		wg.Done()
	}()

	wg.Wait()

	return nil
}
