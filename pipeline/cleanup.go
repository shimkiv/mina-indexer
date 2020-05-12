package pipeline

import (
	"github.com/figment-networks/coda-indexer/config"
	"github.com/figment-networks/coda-indexer/store"
)

// Cleanup pipeline runs the database cleanup steps
type Cleanup struct {
	db  *store.Store
	cfg *config.Config
}

// NewCleanup returns a new cleanup pipeline
func NewCleanup(cfg *config.Config, db *store.Store) *Cleanup {
	return &Cleanup{
		cfg: cfg,
		db:  db,
	}
}

// Execute executes the cleanup pipeline
func (c *Cleanup) Execute() error {
	report, err := c.db.Reports.Last()
	if err != nil {
		// No reports in database
		if err == store.ErrNotFound {
			return nil
		}
		return err
	}

	maxHeight := report.StartHeight - int64(c.cfg.CleanupThreshold)
	if maxHeight <= 0 {
		return nil
	}

	return c.db.Reports.Cleanup(maxHeight)
}
