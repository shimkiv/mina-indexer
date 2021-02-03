package store

import (
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/store/queries"
)

// JobsStore handles operations on jobs
type JobsStore struct {
	baseStore
}

// ByHeight returns all jobs for a given height
func (s JobsStore) ByHeight(height uint64) ([]model.SnarkJob, error) {
	result := []model.SnarkJob{}

	err := s.db.
		Where("height = ?", height).
		Order("id ASC").
		Find(&result).
		Error

	return result, err
}

func (s JobsStore) Import(jobs []model.SnarkJob) error {
	if len(jobs) == 0 {
		return nil
	}

	return bulk.Import(s.db, queries.SnarkJobsImport, len(jobs), func(idx int) bulk.Row {
		j := jobs[idx]

		return bulk.Row{
			j.Height,
			j.Time,
			j.Prover,
			j.Fee,
			j.WorksCount,
			time.Now(),
		}
	})
}
