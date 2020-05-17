package store

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/figment-networks/coda-indexer/model"
)

// Store handles all database operations
type Store struct {
	db *gorm.DB

	Syncables    SyncablesStore
	Reports      ReportsStore
	Blocks       BlocksStore
	Accounts     AccountsStore
	Validators   ValidatorsStore
	Transactions TransactionsStore
	States       StatesStore
	Jobs         JobsStore
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

	// Temporarily disable to let the services start without failing
	// if err := conn.DB().Ping(); err != nil {
	// 	conn.Close()
	// 	return nil, err
	// }

	return &Store{
		db: conn,

		Syncables:    NewSyncablesStore(conn),
		Reports:      NewReportsStore(conn),
		Blocks:       NewBlocksStore(conn),
		Accounts:     NewAccountsStore(conn),
		Validators:   NewValidatorsStore(conn),
		Transactions: NewTransactionsStore(conn),
		States:       NewStatesStore(conn),
		Jobs:         NewJobsStore(conn),
	}, nil
}

func NewSyncablesStore(db *gorm.DB) SyncablesStore {
	return SyncablesStore{scoped(db, model.Report{})}
}

func NewReportsStore(db *gorm.DB) ReportsStore {
	return ReportsStore{scoped(db, model.Report{})}
}

func NewBlocksStore(db *gorm.DB) BlocksStore {
	return BlocksStore{scoped(db, model.Block{})}
}

func NewAccountsStore(db *gorm.DB) AccountsStore {
	return AccountsStore{scoped(db, model.Account{})}
}

func NewValidatorsStore(db *gorm.DB) ValidatorsStore {
	return ValidatorsStore{scoped(db, model.Validator{})}
}

func NewTransactionsStore(db *gorm.DB) TransactionsStore {
	return TransactionsStore{scoped(db, model.Transaction{})}
}

func NewStatesStore(db *gorm.DB) StatesStore {
	return StatesStore{scoped(db, model.State{})}
}

func NewJobsStore(db *gorm.DB) JobsStore {
	return JobsStore{scoped(db, model.Job{})}
}
