package util

import "time"

// HourInterval returns a time interval for an hour
func HourInterval(t time.Time) (time.Time, time.Time) {
	year, month, day := t.UTC().Date()

	start := time.Date(year, month, day, t.Hour(), 0, 0, 0, time.UTC)
	end := time.Date(year, month, day, t.Hour(), 59, 59, 0, time.UTC)

	return start, end
}

// DayInterval returns a time interval for 24h
func DayInterval(t time.Time) (time.Time, time.Time) {
	year, month, day := t.UTC().Date()

	start := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	end := time.Date(year, month, day, 23, 59, 59, 0, time.UTC)

	return start, end
}
