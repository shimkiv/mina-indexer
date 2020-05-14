package store

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/figment-networks/coda-indexer/model"
)

// BlocksStore handles operations on blocks
type BlocksStore struct {
	db *gorm.DB
}

// BlockIndexParams contains the block search params
// TODO: this should probably moved out of store package
type BlockIndexParams struct {
	Creator     string `form:"creator"`
	Order       string `form:"order"`
	OrderColumn string `form:"order_column"`
	Limit       int    `form:"limit"`
}

func (p *BlockIndexParams) setDefaults() {
	if p.Limit <= 0 {
		p.Limit = 25
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
	if p.OrderColumn == "" {
		p.OrderColumn = "id"
	}
	if p.Order == "" {
		p.Order = "DESC"
	}
}

// Create create a new block record
func (s BlocksStore) Create(block *model.Block) error {
	err := s.db.Model(block).Create(block).Error
	return checkErr(err)
}

// Update updates the existing block record
func (s BlocksStore) Update(block *model.Block) error {
	err := s.db.Model(block).Update(block).Error
	return checkErr(err)
}

// CreateIfNotExists creates the block if it does not exist
func (s BlocksStore) CreateIfNotExists(block *model.Block) error {
	_, err := s.FindByHash(block.Hash)
	if err != nil {
		if err == ErrNotFound {
			return s.Create(block)
		}
		return err
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

// Index returns list of blocks
func (s BlocksStore) Index(params BlockIndexParams) ([]model.Block, error) {
	params.setDefaults()

	result := []model.Block{}

	scope := s.db.
		Model(&model.Block{}).
		Order(fmt.Sprintf("%s %s", params.OrderColumn, params.Order)).
		Limit(params.Limit)

	if params.Creator != "" {
		scope = scope.Where("creator = $1", params.Creator)
	}

	err := scope.Find(&result).Error
	return result, checkErr(err)
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
