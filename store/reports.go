package store

import (
	"github.com/jinzhu/gorm"

	"github.com/figment-networks/coda-indexer/model"
)

// ReportsStore handles operations on reports
type ReportsStore struct {
	db *gorm.DB
}

// Create creates a new report record
func (s ReportsStore) Create(report *model.Report) error {
	err := s.db.Model(report).Create(report).Error
	return checkErr(err)
}

// Update updates an existing report record
func (s ReportsStore) Update(report *model.Report) error {
	err := s.db.Model(report).Update(report).Error
	return checkErr(err)
}

// Cleanup removes any reports with a height lower than the provided one
func (s ReportsStore) Cleanup(maxHeight int64) error {
	return s.db.
		Model(&model.Report{}).
		Where("height < ?", maxHeight).
		Error
}

// Last returns the last report
func (s ReportsStore) Last() (*model.Report, error) {
	result := &model.Report{}

	err := s.db.
		Model(result).
		Order("id DESC").
		First(result).Error

	return result, checkErr(err)
}
