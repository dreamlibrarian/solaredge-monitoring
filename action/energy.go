package action

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/dreamlibrarian/solaredge-monitoring/client"
	"github.com/rs/zerolog/log"
)

type EnergyAction struct {
	Action
}

type EnergyConfig struct {
	StartTime time.Time
	EndTime   time.Time
	TimeUnit  string

	DiscoverSites bool
	SiteIDs       []string

	OutputDir string
}

func NewEnergyAction(key string) *EnergyAction {
	return &EnergyAction{
		Action{
			client: client.NewClient(key),
		},
	}
}

// FIXME: Do should return raw objects so I can decide how to format them, and
// maybe do other villainy. this bytestream thing is just bad abstraction.
func (a *EnergyAction) Do(config *EnergyConfig) (map[string][]byte, error) {

	siteIDContentMap := make(map[string][]byte)

	log.Debug().Msg("Getting Energy readings")
	if config.DiscoverSites {
		if len(config.SiteIDs) > 0 {
			return nil, errors.New("cannot set all-sites and specify site-ids")
		}
		log.Debug().Msg("Getting sites from upstream.")
		siteList, err := a.client.GetSiteList()
		if err != nil {
			return nil, err
		}
		log.Debug().Interface("sites", siteList).Msg("got sites")
		for _, site := range siteList {
			config.SiteIDs = append(config.SiteIDs, strconv.FormatInt(site.ID, 10))
		}
	} else if len(config.SiteIDs) == 0 {
		return nil, errors.New("must have at least one site configured, or choose site discovery")
	}

	for _, siteID := range config.SiteIDs {
		usage, err := a.client.GetEnergyUsage(siteID, config.TimeUnit, config.StartTime, config.EndTime)
		if err != nil {
			return nil, err
		}

		usageJSON, err := json.Marshal(usage)
		if err != nil {
			return nil, err
		}

		siteIDContentMap[siteID] = usageJSON

	}
	return siteIDContentMap, nil
}
