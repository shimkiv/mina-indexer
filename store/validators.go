package store

import (
	"time"

	"github.com/figment-networks/coda-indexer/model"
	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/figment-networks/indexing-engine/store/jsonquery"
)

// ValidatorsStore handles operations on validators
type ValidatorsStore struct {
	baseStore
}

func (s ValidatorsStore) Index() ([]byte, error) {
	return jsonquery.MustArray(s.db, sqlValidatorsIndex)
}

// FindAll returns all available validators
func (s ValidatorsStore) FindAll() (result []model.Validator, err error) {
	err = s.db.Order("blocks_created DESC").Find(&result).Error
	return
}

// CreateIfNotExists creates the validator if it does not exist
func (s ValidatorsStore) CreateIfNotExists(validator *model.Validator) error {
	_, err := s.FindByAccount(validator.Account)
	if isNotFound(err) {
		return s.Create(validator)
	}
	return nil
}

// FindByAccount returns a validator record associated with a key
func (s ValidatorsStore) FindByAccount(key string) (*model.Validator, error) {
	result := &model.Validator{}
	err := findBy(s.db, result, "account", key)
	return result, checkErr(err)
}

func (s ValidatorsStore) Import(records []model.Validator) error {
	if len(records) == 0 {
		return nil
	}

	return bulk.Import(s.db, sqlValidatorsImport, len(records), func(idx int) bulk.Row {
		r := records[idx]
		now := time.Now()

		return bulk.Row{
			r.Account,
			r.StartHeight,
			r.StartTime,
			r.LastHeight,
			r.LastTime,
			r.BlocksProposed,
			r.BlocksCreated,
			0,
			0,
			now,
			now,
		}
	})
}

var (
	sqlValidatorsIndex = `
		SELECT
			validators.*,
			accounts.balance AS account_balance,
			accounts.balance_unknown AS account_balance_unknown
		FROM
			validators
		LEFT JOIN accounts
			ON accounts.public_key = validators.account
		ORDER BY
			blocks_created DESC`

	sqlValidatorsImport = `
		INSERT INTO validators (
			account,
			start_height,
			start_time,
			last_height,
			last_time,
			blocks_proposed,
			blocks_created,
			delegated_accounts,
			delegated_balance,
			created_at,
			updated_at
		)
		VALUES
			@values
		ON CONFLICT (account) DO UPDATE
		SET
			last_height        = excluded.last_height,
			last_time          = excluded.last_time,
			blocks_proposed    = excluded.blocks_proposed,
			blocks_created     = COALESCE((SELECT COUNT(1) FROM blocks WHERE creator = excluded.account), 0),
			delegated_accounts = COALESCE((SELECT COUNT(1) FROM accounts WHERE delegate = excluded.account), 0),
			delegated_balance  = COALESCE((SELECT SUM(balance::NUMERIC) FROM accounts WHERE delegate = excluded.account), 0),
			updated_at         = excluded.updated_at`
)
