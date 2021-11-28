package api

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/equipment-data.json
var equipmentDataData []byte

func TestEquipmentDataParse(t *testing.T) {

	var equipmentDataDocument EquipmentDataDocument

	err := json.Unmarshal(equipmentDataData, &equipmentDataDocument)
	if assert.NoError(t, err, "could not parse equipmentDataDocument") {
		assert.NotZero(t, len(equipmentDataDocument.Data.Telemetries))
		for _, telemetry := range equipmentDataDocument.Data.Telemetries {
			if telemetry.Date != nil {
				assert.NotZero(t, telemetry.Date)
			}
		}
	}

}
