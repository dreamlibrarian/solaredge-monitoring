package api

import (
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05"
	dateFormat = "2006-01-02"

	/*
		TimeUnitDay         = "DAY"
		TimeUnitHour        = "HOUR"
		TimeUnitQuarterHour = "QUARTER_OF_AN_HOUR"
	*/
)

func ToDatestamp(t time.Time) string {
	return t.Format(dateFormat)
}

func ToTimestamp(t time.Time) string {
	return t.Format(timeFormat)
}

func parseDate(stamp string) (*time.Time, error) {
	val, err := time.Parse(dateFormat, stamp)
	return &val, err
}

func parseTime(stamp string) (*time.Time, error) {
	val, err := time.Parse(timeFormat, stamp)
	return &val, err
}
