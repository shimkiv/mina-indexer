package store

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/figment-networks/coda-indexer/model"
)

// SyncablesStore handles operations on syncables
type SyncablesStore struct {
	db *gorm.DB
}

// Create creates a new syncable record
func (s SyncablesStore) Create(syncable *model.Syncable) error {
	err := s.db.Model(syncable).Create(syncable).Error
	return checkErr(err)
}

// Update updates an existing syncable record
func (s SyncablesStore) Update(syncable *model.Syncable) error {
	err := s.db.Model(syncable).Update(syncable).Error
	return checkErr(err)
}

// Exists returns true if a syncable of a given kind exists at give height
func (s SyncablesStore) Exists(kind string, height int64) (bool, error) {
	result := &model.Syncable{}

	err := s.db.
		Model(result).
		Where("processed_at IS NOT NULL").
		Where("type = ? AND height = ?", kind, height).
		First(result).
		Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// Count returns the total number of syncables
func (s SyncablesStore) Count(kind string) (int, error) {
	var result int

	err := s.db.
		Model(model.Syncable{}).
		Where("type = ?", kind).
		Count(&result).
		Error

	return result, checkErr(err)
}

// MarkProcessed updates the processed timestamp and saves the changes
func (s SyncablesStore) MarkProcessed(syncable *model.Syncable) error {
	now := time.Now()
	syncable.ProcessedAt = &now

	return s.Update(syncable)
}

// GetMostRecent returns the most recent processed syncable for type
func (s SyncablesStore) GetMostRecent(kind string) (*model.Syncable, error) {
	result := &model.Syncable{}

	err := s.db.
		Model(result).
		Where("processed_at IS NOT NULL").
		Order("height DESC").
		First(result).
		Error

	return result, checkErr(err)
}

// GetMostRecentHeight returns the lowest most recent processed height
func (s SyncablesStore) GetMostRecentHeight() (int64, error) {
	syncable, err := s.GetMostRecent(model.SyncableTypeBlock)
	if err != nil {
		return -1, checkErr(err)
	}
	return syncable.Height, nil
}
