package store

import (
	"fmt"

	"github.com/figment-networks/coda-indexer/model"
)

// BlocksStore handles operations on blocks
type BlocksStore struct {
	baseStore
}

// BlockSearch contains a block search params
type BlockSearch struct {
	Creator string `form:"creator"`

	Order       string `form:"order"`
	OrderColumn string `form:"order_column"`
	Limit       int    `form:"limit"`
}

// CreateIfNotExists creates the block if it does not exist
func (s BlocksStore) CreateIfNotExists(block *model.Block) error {
	_, err := s.FindByHash(block.Hash)
	if isNotFound(err) {
		return s.Create(block)
	}
	return nil
}

// FindBy returns a block for a matching attribute
func (s BlocksStore) FindBy(key string, value interface{}) (*model.Block, error) {
	result := &model.Block{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// FindByID returns a block with matching ID
func (s BlocksStore) FindByID(id int64) (*model.Block, error) {
	return s.FindBy("id", id)
}

// FindByHash returns a block with the matching hash
func (s BlocksStore) FindByHash(hash string) (*model.Block, error) {
	return s.FindBy("hash", hash)
}

// FindByHeight returns a block with the matching height
func (s BlocksStore) FindByHeight(height int64) (*model.Block, error) {
	return s.FindBy("height", height)
}

// Recent returns the most recent block
func (s BlocksStore) Recent() (*model.Block, error) {
	block := &model.Block{}

	err := s.db.
		Order("height DESC").
		First(block).
		Error

	return block, checkErr(err)
}

// Search returns blocks that match search filters
func (s BlocksStore) Search(search BlockSearch) ([]model.Block, error) {
	if search.Limit <= 0 {
		search.Limit = 25
	}
	if search.Limit > 100 {
		search.Limit = 100
	}
	if search.OrderColumn == "" {
		search.OrderColumn = "id"
	}
	if search.Order == "" {
		search.Order = "DESC"
	}

	result := []model.Block{}

	scope := s.db.
		Order(fmt.Sprintf("%s %s", search.OrderColumn, search.Order)).
		Limit(search.Limit)

	if search.Creator != "" {
		scope = scope.Where("creator = ?", search.Creator)
	}

	err := scope.Find(&result).Error
	return result, err
}

// AvgRecentTimes returns recent blocks averages
func (s BlocksStore) AvgRecentTimes(limit int64) (*model.BlockAvgStat, error) {
	res := &model.BlockAvgStat{}

	err := s.db.
		Raw(blockTimesForRecentBlocksQuery, limit).
		Scan(res).
		Error

	return res, checkErr(err)
}

// AvgTimesForInterval returns block stats for a given interval
func (s BlocksStore) AvgTimesForInterval(interval, period string) ([]model.BlockIntervalStat, error) {
	rows, err := s.db.Raw(blockTimesForIntervalQuery, interval, period).Rows()
	if err != nil {
		return nil, checkErr(err)
	}
	defer rows.Close()

	result := []model.BlockIntervalStat{}

	for rows.Next() {
		row := model.BlockIntervalStat{}
		if err := s.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, err
}
