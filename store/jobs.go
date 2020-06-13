package store

import (
	"time"

	"github.com/figment-networks/coda-indexer/model"
	"github.com/figment-networks/indexing-engine/store/bulk"
)

// JobsStore handles operations on jobs
type JobsStore struct {
	baseStore
}

// ByHeight returns all jobs for a given height
func (s JobsStore) ByHeight(height uint64) ([]model.Job, error) {
	result := []model.Job{}

	err := s.db.
		Where("height = ?", height).
		Order("id ASC").
		Take(&result).
		Error

	return result, err
}

func (s JobsStore) Import(jobs []model.Job) error {
	if len(jobs) == 0 {
		return nil
	}

	return bulk.Import(s.db, sqlJobsImport, len(jobs), func(idx int) bulk.Row {
		j := jobs[idx]
		now := time.Now()

		return bulk.Row{
			j.Height,
			j.Time,
			j.Prover,
			j.Fee,
			j.WorksCount,
			now,
			now,
		}
	})
}

var (
	sqlJobsImport = `
		INSERT INTO jobs (
			height,
			time,
			prover,
			fee,
			works_count,
			created_at,
			updated_at
		)
		VALUES @values`
)
