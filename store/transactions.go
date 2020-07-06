package store

import (
	"github.com/figment-networks/coda-indexer/model"
	"github.com/figment-networks/indexing-engine/store/jsonquery"
)

// TransactionsStore handles operations on transactions
type TransactionsStore struct {
	baseStore
}

// CreateIfNotExists creates the transaction if it does not exist
func (s TransactionsStore) CreateIfNotExists(t *model.Transaction) error {
	_, err := s.FindByHash(t.Hash)
	if isNotFound(err) {
		return s.Create(t)
	}
	return nil
}

// FindBy returns transactions by a given key and value
func (s TransactionsStore) FindBy(key string, value interface{}) (*model.Transaction, error) {
	result := &model.Transaction{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// FindByID returns a transaction for a given ID
func (s TransactionsStore) FindByID(id int64) (*model.Transaction, error) {
	return s.FindBy("id", id)
}

// FindByHash returns a transaction for a given hash
func (s TransactionsStore) FindByHash(hash string) (*model.Transaction, error) {
	return s.FindBy("hash", hash)
}

// ByAccount returns a list of transactions sent or received by the account
func (s TransactionsStore) ByAccount(account string) ([]model.Transaction, error) {
	return s.Search(TransactionSearch{Account: account})
}

// ByHeight returns transactions for a given height
func (s TransactionsStore) ByHeight(height uint64) ([]model.Transaction, error) {
	return s.Search(TransactionSearch{Height: height, Limit: 100})
}

// Types returns the list of available transactions types
func (s TransactionsStore) Types() ([]byte, error) {
	return jsonquery.MustObject(s.db, sqlTransactionTypes)
}

// Search returns a list of transactions that matches the filters
func (s TransactionsStore) Search(search TransactionSearch) ([]model.Transaction, error) {
	scope := s.db.
		Order("id DESC").
		Limit(search.Limit)

	if search.BeforeID > 0 {
		scope = scope.Where("id < ?", search.BeforeID)
	}
	if search.AfterID > 0 {
		scope.Where("id > ?", search.AfterID)
	}
	if search.BlockHash != "" {
		scope = scope.Where("block_hash = ?", search.BlockHash)
	}
	if search.Height > 0 {
		scope = scope.Where("height = ?", search.Height)
	}
	if search.Type != "" {
		scope = scope.Where("type = ?", search.Type)
	}
	if search.Account != "" {
		scope = scope.Where("sender = ? OR receiver = ?", search.Account, search.Account)
	} else {
		if search.Sender != "" {
			scope = scope.Where("sender = ?", search.Sender)
		}
		if search.Receiver != "" {
			scope = scope.Where("receiver = ?", search.Receiver)
		}
	}
	if search.Memo != "" {
		scope = scope.Where("memo @@ ?", search.Memo)
	}
	if search.startTime != nil {
		scope = scope.Where("time >= ?", search.startTime)
	}
	if search.endTime != nil {
		scope = scope.Where("time <= ?", search.endTime)
	}

	result := []model.Transaction{}
	err := scope.Find(&result).Error

	return result, err
}

var (
	sqlTransactionTypes = `
		SELECT
  		ARRAY_AGG(e.enumlabel) AS types
		FROM
  		pg_type t
		JOIN pg_enum e
			ON t.oid = e.enumtypid
		JOIN pg_catalog.pg_namespace n
			ON n.oid = t.typnamespace
		WHERE
			t.typname = 'e_tx_type'`
)
