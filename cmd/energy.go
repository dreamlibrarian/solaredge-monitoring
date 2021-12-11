package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/dreamlibrarian/solaredge-monitoring/api"
	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var energyCmd = &cobra.Command{
	Use:   "get-energy",
	Short: "Get energy readings for site",
	RunE:  getEnergy,
}

func getEnergy(cmd *cobra.Command, args []string) error {

	log.Debug().Msg("Getting Energy readings")

	config, err := getEnergyConfig()
	if err != nil {
		return err
	}

	cmd.SilenceUsage = true

	if config.DiscoverSites {
		log.Debug().Msg("Getting sites from upstream.")
		siteList, err := solaredgeClient.GetSiteList()
		if err != nil {
			return err
		}
		log.Debug().Interface("sites", siteList).Msg("got sites")
		for _, site := range siteList {
			config.SiteIDs = append(config.SiteIDs, strconv.FormatInt(site.ID, 10))
		}
	}

	log.Debug().Interface("Config", config).Msg("Got config")

	for _, siteID := range config.SiteIDs {
		usage, err := solaredgeClient.GetEnergyUsage(siteID, config.TimeUnit, config.StartTime, config.EndTime)
		if err != nil {
			return err
		}

		usageJSON, err := json.Marshal(usage)
		if err != nil {
			return err
		}

		filename := fmt.Sprintf("%s/%s_energy.json", config.OutputDir, siteID)
		err = ioutil.WriteFile(filename, usageJSON, 0644)
		if err != nil {
			return fmt.Errorf("unable to write equipment file %s: %w", filename, err)
		}

	}
	return nil
}

type energyConfig struct {
	StartTime time.Time
	EndTime   time.Time
	TimeUnit  string

	SiteIDs []string

	DiscoverSites bool

	OutputDir string
}

func getEnergyConfig() (*energyConfig, error) {
	config := energyConfig{}
	var err, errs error

	startTime := viper.GetString("start-time")
	if startTime == "" {
		config.StartTime = time.Now().Add(-24 * time.Hour)
	} else {
		if config.StartTime, err = parseTime(startTime); err != nil {
			errs = multierror.Append(errs, fmt.Errorf("start time could not be parsed: %w", err))
		}
	}

	endTime := viper.GetString("end-time")
	if endTime == "" {
		config.EndTime = time.Now()
	} else {
		if config.EndTime, err = parseTime(endTime); err != nil {
			errs = multierror.Append(errs, fmt.Errorf("start time could not be parsed: %w", err))
		}
	}

	byHour := viper.GetBool("by-hour")
	byDay := viper.GetBool("by-day")
	byQuarterHour := viper.GetBool("by-quarter-hour")
	var timeSetCount int
	if byHour {
		config.TimeUnit = api.TimeUnitHour
		timeSetCount++
	}
	if byDay {
		config.TimeUnit = api.TimeUnitDay
		timeSetCount++
	}
	if byQuarterHour {
		config.TimeUnit = api.TimeUnitQuarterHour
		timeSetCount++
	}
	if timeSetCount == 0 {
		config.TimeUnit = api.TimeUnitHour
	} else if timeSetCount > 1 {
		errs = multierror.Append(errs, errors.New("may only set one of by-day, by-hour, by-quarter-hour"))
	}

	discoverSites := viper.GetBool("all-sites")
	siteIDs := viper.GetStringSlice("site-id")
	if discoverSites {
		if len(siteIDs) > 0 {
			errs = multierror.Append(errs, errors.New("cannot set all-sites and specify site-ids"))
		} else {
			config.DiscoverSites = true
		}
	} else {
		if len(siteIDs) == 0 {
			errs = multierror.Append(errs, errors.New("must set all-sites or specify at least one site-id"))
		}
		config.SiteIDs = siteIDs
	}

	outputDir := viper.GetString("output-dir")
	outputDirStat, err := os.Stat(outputDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(outputDir, 0755)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("unable to create output directory %s: %w", outputDir, err))
		}
	} else if !outputDirStat.IsDir() {
		errs = multierror.Append(errs, fmt.Errorf("path %s must refer to a directory", outputDir))
	}
	config.OutputDir = outputDir

	return &config, errs
}

func init() {
	RootCmd.AddCommand(telemetryCmd)

	energyCmd.Flags().StringP("start-time", "", "", "Specify the start time for energy - will default to 24 hours ago.")
	energyCmd.Flags().StringP("end-time", "", "", "Specify the start time for energy - will default to now.")
	energyCmd.Flags().BoolP("by-hour", "", false, "Specify hourly samples")
	energyCmd.Flags().BoolP("by-quarter-hour", "", false, "Specify 15-minute samples")
	energyCmd.Flags().BoolP("by-day", "", false, "Specify daily samples")

	energyCmd.Flags().StringSliceP("site-id", "", []string{}, "Specify site IDs; use multiple flags for multiple sites")

	energyCmd.Flags().StringP("output-dir", "", "", "Specify where output files belong")
}
