package store

import (
	"github.com/figment-networks/coda-indexer/model"
)

// ValidatorsStore handles operations on validators
type ValidatorsStore struct {
	baseStore
}

// CreateIfNotExists creates the validator if it does not exist
func (s ValidatorsStore) CreateIfNotExists(validator *model.Validator) error {
	_, err := s.FindByKey(validator.PublicKey)
	if isNotFound(err) {
		return s.Create(validator)
	}
	return nil
}

// FindByKey returns a validator record associated with a key
func (s ValidatorsStore) FindByKey(key string) (*model.Validator, error) {
	result := &model.Validator{}
	err := findBy(s.db, result, "public_key", key)
	return result, checkErr(err)
}
