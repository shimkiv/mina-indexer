package store

import (
	"github.com/figment-networks/coda-indexer/model"
)

// AccountsStore handles operations on accounts
type AccountsStore struct {
	baseStore
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

// FindByID returns an account for the ID
func (s AccountsStore) FindByID(id int64) (*model.Account, error) {
	return s.FindBy("id", id)
}

// FindByPublicKey returns an account for the public key
func (s AccountsStore) FindByPublicKey(key string) (*model.Account, error) {
	return s.FindBy("public_key", key)
}

// ByHeight returns all accounts that were created at a given height
func (s AccountsStore) ByHeight(height int64) ([]model.Account, error) {
	result := []model.Account{}

	err := s.db.
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
		Order("id ASC").
		Find(&result).
		Error

	return result, checkErr(err)
}
