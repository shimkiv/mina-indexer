package util

import (
	"strconv"
	"time"
)

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
