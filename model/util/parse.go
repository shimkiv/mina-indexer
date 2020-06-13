package util

import (
	"strconv"
	"time"
)

// ParseUInt64 returns an UInt64 from a string
func ParseUInt64(input string) (uint64, error) {
	return strconv.ParseUint(input, 10, 64)
}

// ParseInt64 returns an Int64 from a string
func ParseInt64(input string) (int64, error) {
	return strconv.ParseInt(input, 10, 64)
}

// ParseTime returns a timestamp from a string
func ParseTime(input string) (*time.Time, error) {
	msec, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return nil, err
	}
	t := time.Unix(0, msec*1000000)
	return &t, nil
}

// MustInt64 returns an Int64 value without an error
func MustInt64(input string) int64 {
	v, err := ParseInt64(input)
	if err != nil {
		v = 0
	}
	return v
}

// MustUInt64 returns an UInt64 value without an error
func MustUInt64(input string) uint64 {
	v, err := ParseUInt64(input)
	if err != nil {
		v = 0
	}
	return v
}

// MustTime returns a time from a string
func MustTime(input string) time.Time {
	t, err := ParseTime(input)
	if err != nil {
		return time.Time{}
	}
	return *t
}
