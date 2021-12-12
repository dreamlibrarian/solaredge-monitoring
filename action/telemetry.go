package action

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dreamlibrarian/solaredge-monitoring/api"
	"github.com/dreamlibrarian/solaredge-monitoring/client"
	"github.com/rs/zerolog/log"
)

type TelemetryAction struct {
	Action
}

type TelemetryActionConfig struct {
	StartTime time.Time
	EndTime   time.Time
	TimeUnit  string

	SiteIDs       []string
	SerialNumbers []string

	DiscoverSites   bool
	DiscoverSerials bool
}

func NewTelemetryAction(key string) *TelemetryAction {
	return &TelemetryAction{
		Action: Action{
			client: client.NewClient(key),
		},
	}
}

func (t *TelemetryAction) Do(config *TelemetryActionConfig) (map[string]map[string][]api.Telemetry, error) {

	siteIDSerialInventoryMap := make(map[string]map[string][]api.Telemetry)

	log.Debug().Msg("Getting Telemetry")

	log.Debug().Interface("Config", config).Msg("Got config")
	if config.DiscoverSites {
		if len(config.SiteIDs) > 0 {
			return nil, errors.New("cannot set all-sites and specify site-ids")
		}
		log.Debug().Msg("Getting sites from upstream.")
		siteList, err := t.client.GetSiteList()
		if err != nil {
			return nil, err
		}
		log.Debug().Interface("sites", siteList).Msg("got sites")
		for _, site := range siteList {
			config.SiteIDs = append(config.SiteIDs, strconv.FormatInt(site.ID, 10))
		}
	} else if len(config.SiteIDs) == 0 {
		return nil, errors.New("must set all-sites or specify at least one site-id")
	}

	if config.DiscoverSerials && len(config.SerialNumbers) > 0 {
		return nil, errors.New("cannot discover serials and specify serials")
	} else if len(config.SerialNumbers) == 0 {
		return nil, errors.New("must set all-equipment or specify at least one serial")
	}

	for _, siteID := range config.SiteIDs {
		var serialNumbers []string

		siteIDSerialInventoryMap[siteID] = make(map[string][]api.Telemetry)

		if config.DiscoverSerials {
			inventory, err := t.client.GetSiteInventory(siteID)
			if err != nil {
				return nil, fmt.Errorf("unable to get inventory for site %s: %w", siteID, err)
			}
			log.Debug().Interface("inventory", inventory).Msg("got inventory")
			for _, inverter := range inventory.Inverters {
				serialNumbers = append(serialNumbers, inverter.SerialNumber)
			}
			for _, gateway := range inventory.Gateways {
				serialNumbers = append(serialNumbers, gateway.SerialNumber)
			}
			for _, battery := range inventory.Batteries {
				serialNumbers = append(serialNumbers, battery.SerialNumber)
			}
			for _, meter := range inventory.Meters {
				serialNumbers = append(serialNumbers, meter.SerialNumber)
			}
			for _, inverter := range inventory.ThirdPartyInverters {
				serialNumbers = append(serialNumbers, inverter.SerialNumber)
			}
		}

		if len(serialNumbers) == 0 {
			log.Error().Msg("got no serials for site, I'm willing to bet something's wrong")
		}

		for _, serial := range serialNumbers {
			log := log.With().Str("siteid", siteID).Str("serial", serial).Logger()
			equipment, err := t.client.GetTelemetryForEquipment(siteID, serial, config.TimeUnit, config.StartTime, config.EndTime)
			if err != nil {
				return nil, err
			}

			siteIDSerialInventoryMap[siteID][serial] = equipment

			log.Debug().Interface("equipment", equipment).Msg("Got equipment telemetry")
		}

	}

	return siteIDSerialInventoryMap, nil
}
