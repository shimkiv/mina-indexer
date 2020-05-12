package store

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var (
	ErrNotFound = errors.New("record not found")
)

func checkErr(err error) error {
	if gorm.IsRecordNotFoundError(err) {
		return ErrNotFound
	}
	return err
}
