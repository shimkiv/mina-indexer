package store

import (
	"github.com/figment-networks/coda-indexer/model"
)

// TransactionsStore handles operations on transactions
type TransactionsStore struct {
	baseStore
}

// TransactionSearch contains transaction search params
type TransactionSearch struct {
	Height    int64  `form:"height"`
	Type      string `form:"type"`
	BlockHash string `form:"block_hash"`
	Account   string `form:"account"`
	From      string `form:"from"`
	To        string `form:"to"`
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

// Search returns a list of transactions that matches the filters
func (s TransactionsStore) Search(search TransactionSearch) ([]model.Transaction, error) {
	scope := s.db.
		Order("id DESC").
		Limit(100)

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
		scope = scope.Where("sender_key = ? OR recipient_key = ?", search.Account, search.Account)
	} else {
		if search.From != "" {
			scope = scope.Where("sender_key = ?", search.From)
		}
		if search.To != "" {
			scope = scope.Where("recipient_key = ?", search.To)
		}
	}

	result := []model.Transaction{}
	err := scope.Find(&result).Error

	return result, err
}

// ByAccount returns a list of transactions sent or received by the account
func (s TransactionsStore) ByAccount(account string) ([]model.Transaction, error) {
	return s.Search(TransactionSearch{Account: account})
}

// ByHeight returns transactions for a given height
func (s TransactionsStore) ByHeight(height int64) ([]model.Transaction, error) {
	return s.Search(TransactionSearch{Height: height})
}
