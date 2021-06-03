package store

import (
	"fmt"

	"github.com/figment-networks/indexing-engine/store/jsonquery"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/store/queries"
)

// BlocksStore handles operations on blocks
type BlocksStore struct {
	baseStore
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
func (s BlocksStore) FindByHeight(height uint64) (*model.Block, error) {
	return s.FindBy("height", height)
}

// Recent returns the most recent block
func (s BlocksStore) Recent() (*model.Block, error) {
	block := &model.Block{}
	err := s.db.Order("height DESC").Limit(1).Take(block).Error
	return block, checkErr(err)
}

// Search returns blocks that match search filters
func (s BlocksStore) Search(search *BlockSearch) ([]model.Block, error) {
	result := []model.Block{}

	scope := s.db.
		Order(fmt.Sprintf("%s %s", search.Sort, search.Order)).
		Limit(search.Limit)

	if search.MinHeight > 0 {
		scope = scope.Where("height >= ?", search.MinHeight)
	}

	if search.MaxHeight > 0 {
		scope = scope.Where("height <= ?", search.MaxHeight)
	}

	if search.Creator != "" {
		scope = scope.Where("creator = ?", search.Creator)
	}

	return result, scope.Find(&result).Error
}

// AvgTimes returns recent blocks averages
func (s BlocksStore) AvgTimes(limit int64) ([]byte, error) {
	return jsonquery.MustObject(s.db, queries.BlocksTimes, limit)
}

// Stats returns block stats for a given interval
func (s BlocksStore) Stats(period uint, interval string) ([]byte, error) {
	return jsonquery.MustArray(s.db, queries.BlocksStats, period, interval)
}

// FirstBlockOfEpoch returns the first block of epoch
func (s BlocksStore) FirstBlockOfEpoch(epoch string) (*model.Block, error) {
	block := &model.Block{}
	err := s.db.Where("epoch >= ?", epoch).Order("height ASC").Limit(1).Take(block).Error
	return block, checkErr(err)
}
