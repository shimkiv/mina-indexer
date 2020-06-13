package model

import (
	"time"
)

// Model contains the basic model data
type Model struct {
	ID        int64     `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
