package store

import (
	"errors"
)

// BlockSearch contains a block search params
type BlockSearch struct {
	Creator   string `form:"creator"`
	MinHeight uint   `form:"min_height"`
	MaxHeight uint   `form:"max_height"`
	Sort      string `form:"sort"`
	Order     string `form:"order"`
	Limit     uint   `form:"limit"`
}

// Validate performs validation on search parameters
func (search *BlockSearch) Validate() error {
	switch search.Sort {
	case "height":
	case "":
		search.Sort = "height"
	default:
		return errors.New("invalid sort field")
	}

	switch search.Order {
	case "":
		search.Order = "desc"
	case "asc", "desc":
	default:
		return errors.New("invalid sort order")
	}

	if search.Limit == 0 {
		search.Limit = 100
	}
	if search.Limit > 100 {
		return errors.New("max limit is 100")
	}

	return nil
}
