package cli

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/config"
	"github.com/figment-networks/coda-indexer/store"
	"github.com/figment-networks/coda-indexer/worker"
)

func startSyncWorker(wg *sync.WaitGroup, cfg *config.Config, db *store.Store) context.CancelFunc {
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	client := coda.NewDefaultClient(cfg.CodaEndpoint)
	timer := time.NewTimer(cfg.SyncDuration())

	go func() {
		defer func() {
			timer.Stop()
			wg.Done()
		}()

		for {
			select {
			case <-timer.C:
				n, _ := worker.RunSync(cfg, db, client)
				if n > 10 {
					timer.Reset(time.Millisecond * 100)
				} else {
					timer.Reset(cfg.SyncDuration())
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return cancel
}

func startCleanupWorker(wg *sync.WaitGroup, cfg *config.Config, db *store.Store) context.CancelFunc {
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(cfg.CleanupDuration())

	go func() {
		defer func() {
			ticker.Stop()
			wg.Done()
		}()

		for {
			select {
			case <-ticker.C:
				worker.RunCleanup(cfg, db)
			case <-ctx.Done():
				return
			}
		}
	}()

	return cancel
}

func startWorker(cfg *config.Config) error {
	log.Println("using coda endpoint", cfg.CodaEndpoint)
	log.Println("sync will run every", cfg.SyncInterval)
	log.Println("cleanup will run every", cfg.CleanupInterval)

	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	wg := &sync.WaitGroup{}

	cancelSync := startSyncWorker(wg, cfg, db)
	cancelCleanup := startCleanupWorker(wg, cfg, db)

	s := <-initSignals()

	log.Println("received signal", s)
	cancelSync()
	cancelCleanup()

	wg.Wait()
	return nil
}
