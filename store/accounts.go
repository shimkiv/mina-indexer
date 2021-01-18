package store

import (
	"time"

	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/mina-indexer/model"
)

// AccountsStore handles operations on accounts
type AccountsStore struct {
	baseStore
}

func (s AccountsStore) Count() (int, error) {
	var n int
	err := s.db.Table("accounts").Count(&n).Error
	return n, err
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

// AllByDelegator returns all accounts delegated to another account
func (s AccountsStore) AllByDelegator(account string) ([]model.Account, error) {
	result := []model.Account{}
	err := s.db.
		Where("delegate = ?", account).
		Find(&result).
		Error
	return result, checkErr(err)
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

func (s AccountsStore) Import(records []model.Account) error {
	if len(records) == 0 {
		return nil
	}

	return bulk.Import(s.db, sqlAccountsImport, len(records), func(idx int) bulk.Row {
		acc := records[idx]
		now := time.Now()

		return bulk.Row{
			acc.PublicKey,
			acc.Delegate,
			acc.Balance,
			acc.BalanceUnknown,
			acc.Nonce,
			acc.StartHeight,
			acc.StartTime,
			acc.LastHeight,
			acc.LastTime,
			now,
			now,
		}
	})
}

var (
	sqlAccountsImport = `
		INSERT INTO accounts (
			public_key,
			delegate,
			balance,
			balance_unknown,
			nonce,
			start_height,
			start_time,
			last_height,
			last_time,
			created_at,
			updated_at
		)
		VALUES @values
		ON CONFLICT (public_key) DO UPDATE
		SET
			delegate        = excluded.delegate,
			balance         = excluded.balance,
			balance_unknown = excluded.balance_unknown,
			nonce           = excluded.nonce,
			last_height     = excluded.last_height,
			last_time       = excluded.last_time,
			updated_at      = excluded.updated_at`
)
