package indexing

import (
	"github.com/figment-networks/coda-indexer/store"
	log "github.com/sirupsen/logrus"
)

// Finalize generates summary records
func Finalize(db *store.Store, data *Data) error {
	ts := data.Block.Time
	buckets := []string{store.BucketHour, store.BucketDay}

	for _, bucket := range buckets {
		log.WithField("bucket", bucket).Debug("creating chain stats")
		if err := db.Stats.CreateChainStats(bucket, ts); err != nil {
			return err
		}

		log.WithField("bucket", bucket).Debug("creating transaction stats")
		if err := db.Stats.CreateTransactionsStats(bucket, ts); err != nil {
			return err
		}
	}

	return nil
}
