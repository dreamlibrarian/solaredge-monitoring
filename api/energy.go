package api

import (
	"encoding/json"
	"fmt"
	"time"
)

type Energy struct {
	TimeUnit   string  `json:"timeUnit"`
	Unit       string  `json:"unit"`
	MeasuredBy string  `json:"measuredBy"`
	Values     []Value `json:"values"`
}

type Value struct {
	Date  time.Time `json:"date"`
	Value *int64    `json:"value"`
}

// UnmarshalJSON for Value has to do the extra work to parse the Date into time.Time
func (v *Value) UnmarshalJSON(data []byte) error {
	interimData := struct {
		Date  string `json:"date"`
		Value *int64 `json:"value"`
	}{}

	err := json.Unmarshal(data, &interimData)
	if err != nil {
		return fmt.Errorf("unable to parse '%s' as Value : %w", string(data), err)
	}

	timeStamp, err := parseTime(interimData.Date)
	if err != nil {
		return fmt.Errorf("unable to parse %s as date format %s: %w", interimData.Date, dateFormat, err)
	}

	*v = Value{Date: *timeStamp, Value: interimData.Value}

	return nil
}
