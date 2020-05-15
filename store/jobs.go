package store

import "github.com/figment-networks/coda-indexer/model"

// JobsStore handles operations on jobs
type JobsStore struct {
	baseStore
}

// ByHeight returns all jobs for a given height
func (s JobsStore) ByHeight(height int64) ([]model.Job, error) {
	result := []model.Job{}

	err := s.db.
		Where("height = ?", height).
		Order("id ASC").
		Error

	return result, err
}
