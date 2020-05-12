package store

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Store handles all database operations
type Store struct {
	db *gorm.DB

	Syncables    SyncablesStore
	Reports      ReportsStore
	Blocks       BlocksStore
	Accounts     AccountsStore
	Transactions TransactionsStore
	States       StatesStore
}

// Test checks the connection status
func (s *Store) Test() error {
	return s.db.DB().Ping()
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// SetDebugMode enabled detailed query logging
func (s *Store) SetDebugMode(enabled bool) {
	s.db.LogMode(enabled)
}

// New returns a new store from the connection string
func New(connStr string) (*Store, error) {
	conn, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := conn.DB().Ping(); err != nil {
		conn.Close()
		return nil, err
	}

	return &Store{
		db: conn,

		Syncables:    SyncablesStore{conn},
		Reports:      ReportsStore{conn},
		Blocks:       BlocksStore{conn},
		Accounts:     AccountsStore{conn},
		Transactions: TransactionsStore{conn},
		States:       StatesStore{conn},
	}, nil
}

func findBy(db *gorm.DB, dst interface{}, key string, value interface{}) error {
	return db.
		Model(dst).
		Where(fmt.Sprintf("%s = ?", key), value).
		First(dst).
		Error
}
