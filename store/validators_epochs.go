package store

import (
	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/store/queries"
)

// ValidatorEpochStore handles operations on validator epochs
type ValidatorsEpochsStore struct {
	baseStore
}

// ImportValidatorEpochs creates or updates validator records in bulk
func (s ValidatorsEpochsStore) Import(records []model.ValidatorEpoch) error {
	if len(records) == 0 {
		return nil
	}

	return bulk.Import(s.db, queries.ValidatorsEpochsImport, len(records), func(idx int) bulk.Row {
		r := records[idx]
		return bulk.Row{
			r.AccountId,
			r.Epoch,
			r.ValidatorFee,
		}
	})
}

// GetValidatorEpochs fetches validator epochs
func (s *ValidatorsEpochsStore) GetValidatorEpochs(epoch string, accountId string) ([]model.ValidatorEpoch, error) {
	var res []model.ValidatorEpoch

	scope := s.db

	if accountId != "" {
		scope = scope.Where("accountId = ?", accountId)
	}
	if epoch != "" {
		scope = scope.Where("epoch = ?", epoch)
	}

	return res, scope.Find(&res).Error
}
