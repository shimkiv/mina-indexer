package types

import (
	"database/sql/driver"
	"errors"
	"math/big"
)

var (
	errInvalidFloat = errors.New("invalid float")
)

type Float struct {
	*big.Float
}

// NewFloat returns a new float from the given string
func NewFloat(src string) Float {
	if src == "" {
		src = "0"
	}

	n := new(big.Float)
	n.SetString(src)

	return Float{Float: n}
}

// NewFloat64Float returns a new float for the given float64 value
func NewFloat64Float(val float64) Float {
	n := big.NewFloat(val)
	return Float{Float: n}
}

// Value returns a serialized value
func (a Float) Value() (driver.Value, error) {
	if a.Float != nil {
		return a.Float.String(), nil
	}
	return nil, nil
}

func (a Float) String() string {
	if a.Float == nil {
		return ""
	}
	return a.Float.String()
}

// Add adds two numbers
func (a Float) Add(b Float) Float {
	n := new(big.Float)
	n = n.Add(a.Float, b.Float)
	return Float{n}
}

// Sub substitutes a given float amount from the current one
func (a Float) Sub(b Float) Float {
	n := new(big.Float)
	n = n.Sub(a.Float, b.Float)
	return Float{n}
}

// Mul multiplies two numbers
func (a Float) Mul(b Float) Float {
	n := new(big.Float)
	n = n.Mul(a.Float, b.Float)
	return Float{n}
}

// Quo divides two numbers
func (a Float) Quo(b Float) Float {
	n := new(big.Float)
	n = n.Quo(a.Float, b.Float)
	return Float{n}
}

// Scan assigns the value from interface
func (a *Float) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	n := new(big.Float)

	switch v := value.(type) {
	case float64:
		a.Float = n.SetFloat64(v)
	case string:
		n, ok := n.SetString(v)
		if !ok {
			return errInvalidFloat
		}
		a.Float = n
	case []byte:
		n, ok := n.SetString(string(value.([]byte)))
		if !ok {
			return errInvalidFloat
		}
		a.Float = n
	}

	return nil
}
