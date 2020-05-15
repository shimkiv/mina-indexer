package store

import (
	"github.com/figment-networks/coda-indexer/model"
)

// StatesStore handles operations on states
type StatesStore struct {
	baseStore
}

// CreateIfNotExists creates a new state record unless it exists
func (s StatesStore) CreateIfNotExists(state *model.State) error {
	_, err := s.FindByHeight(state.Height)
	if isNotFound(err) {
		return s.Create(state)
	}
	return nil
}

// FindByHeight returns a state record for a given height
func (s StatesStore) FindByHeight(height int64) (*model.State, error) {
	result := &model.State{}

	err := s.db.
		Where("height = ?", height).
		First(result).
		Error

	return result, checkErr(err)
}
