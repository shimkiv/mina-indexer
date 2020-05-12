package store

import (
	"github.com/jinzhu/gorm"

	"github.com/figment-networks/coda-indexer/model"
)

// AccountsStore handles operations on accounts
type AccountsStore struct {
	db *gorm.DB
}

// Create creates a new account record
func (s AccountsStore) Create(acc *model.Account) error {
	err := s.db.Model(acc).Create(acc).Error
	return checkErr(err)
}

// Update updates the existing account record
func (s AccountsStore) Update(acc *model.Account) error {
	err := s.db.Model(acc).Update(acc).Error
	return checkErr(err)
}

// CreateOrUpdate creates a new account or updates an existing one
func (s AccountsStore) CreateOrUpdate(acc *model.Account) error {
	existing, err := s.FindByPublicKey(acc.PublicKey)
	if err != nil {
		if err == ErrNotFound {
			return s.Create(acc)
		}
		return err
	}

	existing.Balance = acc.Balance
	existing.Nonce = acc.Nonce

	return s.Update(existing)
}

// FindBy returns an account for a matching attribute
func (s AccountsStore) FindBy(key string, value interface{}) (*model.Account, error) {
	result := &model.Account{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// FindByPublicKey returns an account for the public key
func (s AccountsStore) FindByPublicKey(key string) (*model.Account, error) {
	return s.FindBy("public_key", key)
}

// ListByHeight returns all accounts that were created at a given height
func (s AccountsStore) ListByHeight(height int64) ([]model.Account, error) {
	result := []model.Account{}

	err := s.db.
		Model(&model.Account{}).
		Where("start_height <= ?", height).
		Order("id DESC").
		Find(&result).
		Error

	return result, checkErr(err)
}

// All returns all accounts
func (s AccountsStore) All() ([]model.Account, error) {
	result := []model.Account{}

	err := s.db.
		Model(&model.Account{}).
		Order("id ASC").
		Find(&result).
		Error

	return result, checkErr(err)
}
