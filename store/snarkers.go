package store

import (
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/indexing-engine/store/jsonquery"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/store/queries"
)

type SnarkersStore struct {
	baseStore
}

func (s SnarkersStore) All() ([]model.Snarker, error) {
	result := []model.Snarker{}
	err := s.db.
		Model(&model.Snarker{}).
		Order("jobs_count DESC").
		Find(&result).
		Error
	return result, checkErr(err)
}

// FindSnarker returns snarker for a given account
func (s SnarkersStore) FindSnarker(account string) (*model.Snarker, error) {
	result := &model.Snarker{}
	err := findBy(s.db, result, "account", account)
	return result, checkErr(err)
}

// SnarkerInfoFromCanonicalBlocks returns snarker info from canonical blocks
func (s SnarkersStore) SnarkerInfoFromCanonicalBlocks(account string, start, end uint64) ([]byte, error) {
	return jsonquery.MustObject(s.db, queries.SnarkerInfoFromCanonicalBlocks, account, start, end)
}

func (s SnarkersStore) Import(records []model.Snarker) error {
	if len(records) == 0 {
		return nil
	}

	now := time.Now()

	return bulk.Import(s.db, queries.SnarkersImport, len(records), func(idx int) bulk.Row {
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
