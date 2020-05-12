package store

import (
	"github.com/figment-networks/coda-indexer/model"
	"github.com/jinzhu/gorm"
)

// TransactionsStore handles operations on transactions
type TransactionsStore struct {
	db *gorm.DB
}

// Create creates a new transaction record
func (s TransactionsStore) Create(t *model.Transaction) error {
	err := s.db.Model(t).Create(t).Error
	return checkErr(err)
}

// CreateIfNotExist creates the transaction if it does not exist
func (s TransactionsStore) CreateIfNotExist(t *model.Transaction) error {
	_, err := s.FindByHash(t.Hash)
	if err != nil {
		if err == ErrNotFound {
			return s.Create(t)
		}
		return err
	}
	return nil
}

// Update updates the existing transaction record
func (s TransactionsStore) Update(t *model.Transaction) error {
	err := s.db.Model(t).Update(t).Error
	return checkErr(err)
}

// FindBy returns transactions by a given key and value
func (s TransactionsStore) FindBy(key string, value interface{}) (*model.Transaction, error) {
	result := &model.Transaction{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// FindByHash returns a transaction for a given hash
func (s TransactionsStore) FindByHash(hash string) (*model.Transaction, error) {
	return s.FindBy("hash", hash)
}

// ListByHeight returns transactions for a given height
func (s TransactionsStore) ListByHeight(height int64) ([]model.Transaction, error) {
	result := []model.Transaction{}
	scope := s.db.Model(&model.Transaction{}).Order("id ASC")

	if height > 0 {
		scope = scope.Where("height = ?", height)
	}

	err := scope.Find(&result).Error
	return result, err
}
