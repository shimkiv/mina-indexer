package store

import (
	"github.com/jinzhu/gorm"

	"github.com/figment-networks/coda-indexer/model"
)

// StatesStore handles operations on states
type StatesStore struct {
	db *gorm.DB
}

// GetByHeight returns a state record for a given height
func (s StatesStore) GetByHeight(height int64) (*model.State, error) {
	result := &model.State{}

	err := s.db.
		Model(result).
		Where("height = ?", height).
		First(result).
		Error

	return result, checkErr(err)
}

// Create creates a new state record
func (s StatesStore) Create(state *model.State) error {
	return checkErr(s.db.Model(state).Create(state).Error)
}

// CreateIfNotExists creates a new state record unless it exists
func (s StatesStore) CreateIfNotExists(state *model.State) error {
	_, err := s.GetByHeight(state.Height)
	if err == ErrNotFound {
		return s.Create(state)
	}
	return nil
}
