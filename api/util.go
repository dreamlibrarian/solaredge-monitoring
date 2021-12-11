package api

import (
	"time"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
	DateFormat = "2006-01-02"

	TimeUnitDay         = "DAY"
	TimeUnitHour        = "HOUR"
	TimeUnitQuarterHour = "QUARTER_OF_AN_HOUR"
)

func ToDatestamp(t time.Time) string {
	return t.Format(DateFormat)
}

func ToTimestamp(t time.Time) string {
	return t.Format(TimeFormat)
}

func parseDate(stamp string) (*time.Time, error) {
	val, err := time.Parse(DateFormat, stamp)
	return &val, err
}

func parseTime(stamp string) (*time.Time, error) {
	val, err := time.Parse(TimeFormat, stamp)
	return &val, err
}
