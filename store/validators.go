package store

import (
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/indexing-engine/store/jsonquery"

	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/store/queries"
)

// ValidatorsStore handles operations on validators
type ValidatorsStore struct {
	baseStore
}

func (s ValidatorsStore) Index() ([]byte, error) {
	return jsonquery.MustArray(s.db, queries.ValidatorsIndex)
}

// FindAll returns all available validators
func (s ValidatorsStore) FindAll() (result []model.Validator, err error) {
	err = s.db.Order("blocks_created DESC").Find(&result).Error
	return
}

// FindByPublicKey returns a validator record associated with a key
func (s ValidatorsStore) FindByPublicKey(key string) (*model.Validator, error) {
	result := &model.Validator{}
	err := findBy(s.db, result, "public_key", key)
	return result, checkErr(err)
}

func (s ValidatorsStore) UpdateStaking() error {
	return s.db.Exec(queries.ValidatorsUpdateStaking).Error
}

// UpdateIdentity updates the identity name of the validator
func (s ValidatorsStore) UpdateIdentity(key string, name string) error {
	return s.db.Exec(
		"UPDATE validators SET identity_name = ? WHERE public_key = ?",
		name, key,
	).Error
}

// UpdateFee updates the validator fee
func (s ValidatorsStore) UpdateFee(key string, value float64) error {
	return s.db.Exec(queries.ValidatorEpochsCreate, key, value).Error
}

// Import creates or updates validator records in bulk
func (s ValidatorsStore) Import(records []model.Validator) error {
	if len(records) == 0 {
		return nil
	}

	return bulk.Import(s.db, queries.ValidatorsImport, len(records), func(idx int) bulk.Row {
		r := records[idx]
		now := time.Now()

		return bulk.Row{
			r.PublicKey,
			r.StartHeight,
			r.StartTime,
			r.LastHeight,
			r.LastTime,
			r.BlocksProposed,
			r.BlocksCreated,
			now,
			now,
		}
	})
}
