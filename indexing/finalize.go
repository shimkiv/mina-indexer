package indexing

import (
	"github.com/figment-networks/coda-indexer/store"
	log "github.com/sirupsen/logrus"
)

// Finalize generates summary records
func Finalize(db *store.Store, data *Data) error {
	ts := data.Block.Time

	// Create hourly stats rollup
	log.Debug("creating hourly stats")
	if err := db.Stats.CreateChainStats(store.BucketHour, ts); err != nil {
		return err
	}

	log.Debug("creating daily stats")
	if err := db.Stats.CreateChainStats(store.BucketDay, ts); err != nil {
		return err
	}

	return nil
}
