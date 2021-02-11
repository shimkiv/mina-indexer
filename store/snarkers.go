package store

import (
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/store/queries"
)

type SnarkersStore struct {
	baseStore
}

func (s SnarkersStore) All() ([]model.Snarker, error) {
	result := []model.Snarker{}
	err := s.db.Find(&result).Error
	return result, checkErr(err)
}

func (s SnarkersStore) Import(records []model.Snarker) error {
	if len(records) == 0 {
		return nil
	}

	now := time.Now()

	return bulk.Import(s.db, queries.SnarkersImport, len(records), func(idx int) bulk.Row {
		r := records[idx]

		return bulk.Row{
			r.PublicKey,
			r.Fee,
			r.JobsCount,
			r.WorksCount,
			r.StartHeight,
			r.StartTime,
			r.LastHeight,
			r.LastTime,
			now, now,
		}
	})
}
