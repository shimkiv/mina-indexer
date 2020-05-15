package store

import (
	"github.com/figment-networks/coda-indexer/model"
)

// ReportsStore handles operations on reports
type ReportsStore struct {
	baseStore
}

// Cleanup removes any reports with a height lower than the provided one
func (s ReportsStore) Cleanup(maxHeight int64) error {
	return s.db.Delete(s.model, "height < ?", maxHeight).Error
}

// Last returns the last report
func (s ReportsStore) Last() (*model.Report, error) {
	result := &model.Report{}

	err := s.db.
		Order("id DESC").
		First(result).Error

	return result, checkErr(err)
}
