package store

import (
	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/store/queries"
)

// ValidatorsEpochsStore handles operations on validator epochs
type ValidatorsEpochsStore struct {
	baseStore
}

// Import creates or updates validator records in bulk
func (s ValidatorsEpochsStore) Import(records []model.ValidatorEpoch) error {
	if len(records) == 0 {
		return nil
	}

	return bulk.Import(s.db, queries.ValidatorsEpochsImport, len(records), func(idx int) bulk.Row {
		r := records[idx]
		return bulk.Row{
			r.AccountAddress,
			r.Epoch,
			r.ValidatorFee,
		}
	})
}

// GetValidatorEpochs fetches validator epochs
func (s *ValidatorsEpochsStore) GetValidatorEpochs(epoch string, accountAddress string) ([]model.ValidatorEpoch, error) {
	var res []model.ValidatorEpoch

	scope := s.db

	if accountAddress != "" {
		scope = scope.Where("account_address = ?", accountAddress)
	}
	if epoch != "" {
		scope = scope.Where("epoch = ?", epoch)
	}

	return res, scope.Find(&res).Error
}

// GetLastValidatorEpoch returns the most recent validator epoch record
func (s *ValidatorsEpochsStore) GetLastValidatorEpoch(address string) (*model.ValidatorEpoch, error) {
	result := &model.ValidatorEpoch{}

	err := s.db.
		Model(result).
		Where("account_address = ?", address).
		Order("epoch DESC").
		Limit(1).
		Take(result).
		Error

	return result, checkErr(err)
}
