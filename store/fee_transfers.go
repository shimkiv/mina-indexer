package store

import (
	"time"

	"github.com/figment-networks/coda-indexer/model"
	"github.com/figment-networks/indexing-engine/store/bulk"
)

type FeeTransfersStore struct {
	baseStore
}

func (s FeeTransfersStore) FindByHeight(height uint64) (result []model.FeeTransfer, err error) {
	err = s.db.
		Where("height = ?", height).
		Order("id ASC").
		Take(&result).
		Error
	return
}

func (s FeeTransfersStore) Import(records []model.FeeTransfer) error {
	if len(records) == 0 {
		return nil
	}

	now := time.Now()

	return bulk.Import(s.db, sqlFeeTransfersImport, len(records), func(idx int) bulk.Row {
		return bulk.Row{
			records[idx].Height,
			records[idx].Time,
			records[idx].Recipient,
			records[idx].Amount,
			now,
			now,
		}
	})
}

var (
	sqlFeeTransfersImport = `
		INSERT INTO fee_transfers (
			height,
			time,
			recipient,
			amount,
			created_at,
			updated_at
		)
		VALUES @values`
)
