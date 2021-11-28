package api

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/energy-day-periods.json
var energyDayPeriodTestData []byte

//go:embed testdata/energy-hour-periods.json
var energyHourPeriodTestData []byte

//go:embed testdata/energy-quarter-hour-periods.json
var energyQuarterHourPeriodTestData []byte

func TestEnergyParse(t *testing.T) {
	var err error

	var dayEnergy Energy
	err = json.Unmarshal(energyDayPeriodTestData, &dayEnergy)
	assert.NoErrorf(t, err, "Could not parse day file as Energy")

	var hourEnergy Energy
	err = json.Unmarshal(energyHourPeriodTestData, &hourEnergy)
	assert.NoErrorf(t, err, "Could not parse hour file as Energy")

	var quarterHourEnergy Energy
	err = json.Unmarshal(energyQuarterHourPeriodTestData, &quarterHourEnergy)
	assert.NoErrorf(t, err, "Could not parse quarterHour file as Energy")

}
