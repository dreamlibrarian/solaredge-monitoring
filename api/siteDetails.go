package api

import (
	"encoding/json"
	"fmt"
	"time"
)

type SiteDetailsDocument struct {
	Details SiteDetails `json:"details"`
}

type SiteListDocument struct {
	Sites struct {
		Count int64         `json:"count"`
		Sites []SiteDetails `json:"site"`
	}
}

type SiteDetails struct {
	ID               int64             `json:"id"`
	Name             string            `json:"name"`
	AccountID        int64             `json:"accountId"`
	Status           string            `json:"status"`
	PeakPower        float64           `json:"peakPower"`
	LastUpdateTime   *time.Time        `json:"lastUpdateTime"`
	InstallationDate *time.Time        `json:"installationDate"`
	PTODate          *time.Time        `json:"ptoDate"`
	Notes            string            `json:"notes"`
	Type             string            `json:"type"`
	Location         Location          `json:"location"`
	PrimaryModule    PrimaryModule     `json:"primaryModule"`
	URIs             map[string]string `json:"uris"`
	PublicSettings   PublicSettings    `json:"publicSettings"`
}

func (s *SiteDetails) UnmarshalJSON(data []byte) error {
	type SiteDetailsAlias SiteDetails
	interimData := struct {
		SiteDetailsAlias
		LastUpdateTime   *string `json:"lastUpdateTime"`
		InstallationDate *string `json:"installationDate"`
		PTODate          *string `json:"ptoDate"`
	}{}
	err := json.Unmarshal(data, &interimData)
	if err != nil {
		return fmt.Errorf("unable to unmarshal data %s to interim data structure: %w", data, err)
	}

	*s = SiteDetails(interimData.SiteDetailsAlias)

	if interimData.LastUpdateTime != nil {
		s.LastUpdateTime, err = parseDate(*interimData.LastUpdateTime)
		if err != nil {
			return fmt.Errorf("unable to parse lastUpdateTime %s with format %s: %w", *interimData.LastUpdateTime, timeFormat, err)
		}
	}

	if interimData.InstallationDate != nil {
		s.InstallationDate, err = parseDate(*interimData.InstallationDate)
		if err != nil {
			return fmt.Errorf("unable to parse installationDate %s with format %s: %w", *interimData.InstallationDate, timeFormat, err)
		}
	}

	if interimData.PTODate != nil {
		s.PTODate, err = parseTime(*interimData.PTODate)
		if err != nil {
			return fmt.Errorf("unable to parse PTODate %s with format %s: %w", *interimData.PTODate, timeFormat, err)
		}
	}

	return nil
}

type Location struct {
	Country     string `json:"country"`
	State       string `json:"state"`
	City        string `json:"city"`
	Address     string `json:"address"`
	Address2    string `json:"address2"`
	Zip         string `json:"zip"`
	TimeZone    string `json:"timeZone"`
	CountryCode string `json:"countryCode"`
	StateCode   string `json:"stateCode"`
}

type PrimaryModule struct {
	ManufacturerName       string  `json:"manufacturerName"`
	ModelName              string  `json:"modelName"`
	MaximumPower           float64 `json:"maximumPower"`
	TemperatureCoefficient float64 `json:"temperatureCoef"`
}

type PublicSettings struct {
	IsPublic bool `json:"isPublic"`
}
