package action

import (
	"errors"
	"strconv"
	"time"

	"github.com/dreamlibrarian/solaredge-monitoring/api"
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
}

func NewEnergyAction(key string) *EnergyAction {
	return &EnergyAction{
		Action{
			client: client.NewClient(key),
		},
	}
}

func (a *EnergyAction) Do(config *EnergyConfig) (map[string]*api.Energy, error) {

	siteIDContentMap := make(map[string]*api.Energy)

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

		siteIDContentMap[siteID] = usage
	}
	return siteIDContentMap, nil
}
