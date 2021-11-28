package api

import (
	"encoding/json"
	"fmt"
	"time"
)

type EquipmentDataDocument struct{
	Data EquipmentData `json:"data"`
}

type EquipmentData struct {
	Count       int         `json:"count"`
	Telemetries []Telemetry `json:"telemetries"`
}

type Telemetry struct {
	Date                  *time.Time `json:"date"`
	TotalActivePower      float64   `json:"totalActivePower"`
	DCVoltage             float64   `json:"dcVoltage"`
	GroundFaultResistance int64     `json:"groundFaultResistance"`
	PowerLimit            int64     `json:"powerLimit"`
	TotalEnergy           int64     `json:"totalEnergy"`
	Temperature           float64   `json:"temperature"`
	InverterMode          string    `json:"inverterMode"`
	OperationMode         int64     `json:"operationMode"`
	L1Data                L1Data    `json:"L1Data"`
}

func (t *Telemetry) UnmarshalJSON(data []byte) error {
	type TelemetryAlias Telemetry

	interimTelemetry := struct {
		TelemetryAlias
		Date string `json:"date"`
	}{}

	err := json.Unmarshal(data, &interimTelemetry)
	if err != nil {
		return fmt.Errorf("unable to parse telemetry to interim format: %w", err)
	}

	*t = Telemetry(interimTelemetry.TelemetryAlias)

	t.Date, err = parseTime(interimTelemetry.Date)

	return nil
}

type L1Data struct {
	ACCurrent     float64 `json:"acCurrent"`
	ACVoltage     float64 `json:"ACVoltage"`
	ACFrequency   float64 `json:"ACFrequency"`
	ApparentPower float64 `json:"apparentPower"`
	ActivePower   float64 `json:"activePower"`
	ReactivePower float64 `json:"reactivePower"`
	CosPhi        int64   `json:"cosPhi"`
}
