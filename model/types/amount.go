package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
)

var (
	errInvalidAmount = errors.New("invalid amount")

	zero = new(big.Int)
)

// Amount represense a NEAR yocto
type Amount struct {
	*big.Int
}

// NewAmount returns a new amount from the given string
func NewAmount(src string) Amount {
	if src == "" {
		src = "0"
	}

	n := new(big.Int)
	n.SetString(src, 10)

	return Amount{Int: n}
}

// NewInt64Amount returns a new amount for the given int64 value
func NewInt64Amount(val int64) Amount {
	n := big.NewInt(val)
	return Amount{Int: n}
}

func NewFloatAmount(val string) Amount {
	var n float64
	fmt.Sscanf(val, "%f", &n)

	return NewInt64Amount(int64(n * 1000000000))
}

// MarshalJSON returns a JSON representation of amount
func (a Amount) MarshalJSON() ([]byte, error) {
	if a.Int == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(a.Int.String())
}

// Value returns a serialized value
func (a Amount) Value() (driver.Value, error) {
	if a.Int != nil {
		return a.Int.String(), nil
	}
	return nil, nil
}

func (a Amount) String() string {
	if a.Int == nil {
		return ""
	}
	return a.Int.String()
}

// Compare compares two amounts
func (a Amount) Compare(b Amount) int {
	return a.Cmp(b.Int)
}

// Add adds two numbers
func (a Amount) Add(b Amount) Amount {
	n := new(big.Int)
	n = n.Add(a.Int, b.Int)
	return Amount{n}
}

// Sub substitutes a given amount from the current one
func (a Amount) Sub(b Amount) Amount {
	n := new(big.Int)
	n = n.Sub(a.Int, b.Int)
	return Amount{n}
}

// Mul multiplies two numbers
func (a Amount) Mul(b Amount) Amount {
	n := new(big.Int)
	n = n.Mul(a.Int, b.Int)
	return Amount{n}
}

// Quo divides two numbers
func (a Amount) Quo(b Amount) Amount {
	n := new(big.Int)
	n = n.Quo(a.Int, b.Int)
	return Amount{n}
}

func (a Amount) PercentOf(b Amount) float64 {
	if b.Int.Cmp(zero) == 0 {
		return float64(0.0)
	}

	x := new(big.Float).SetInt(a.Int)
	x = x.Mul(x, big.NewFloat(100.0))
	y := new(big.Float).SetInt(b.Int)

	result, _ := new(big.Float).Quo(x, y).Float64()
	return result
}

// Scan assigns the value from interface
func (a *Amount) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	n := new(big.Int)

	switch v := value.(type) {
	case int64:
		a.Int = n.SetInt64(v)
	case string:
		n, ok := n.SetString(v, 10)
		if !ok {
			return errInvalidAmount
		}
		a.Int = n
	case []byte:
		n, ok := n.SetString(string(value.([]byte)), 10)
		if !ok {
			return errInvalidAmount
		}
		a.Int = n
	}

	return nil
}
