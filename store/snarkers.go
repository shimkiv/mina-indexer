package store

import (
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/mina-indexer/model"
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

	return bulk.Import(s.db, sqlSnarkersImport, len(records), func(idx int) bulk.Row {
		r := records[idx]

		return bulk.Row{
			r.Account,
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

var (
	sqlSnarkersImport = `
		INSERT INTO snarkers (
			account,
			fee,
			jobs_count,
			works_count,
			start_height,
			start_time,
			last_height,
			last_time,
			created_at,
			updated_at
		)
		VALUES @values
		ON CONFLICT (account) DO UPDATE
		SET
			fee         = excluded.fee,
			jobs_count  = snarkers.jobs_count + excluded.jobs_count,
			works_count = snarkers.works_count + excluded.works_count,
			last_height = excluded.last_height,
			last_time   = excluded.last_time,
			updated_at  = excluded.updated_at`
)
