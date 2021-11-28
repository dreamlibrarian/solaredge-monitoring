package api

import (
	"time"
)

const timeFormat = "2006-01-02 03:04:05"
const dateFormat = "2006-01-02"

func parseDate(stamp string) (*time.Time, error) {
	val, err := time.Parse(dateFormat, stamp)
	return &val, err
}

func parseTime(stamp string) (*time.Time, error) {
	val, err := time.Parse(timeFormat, stamp)
	return &val, err
}
