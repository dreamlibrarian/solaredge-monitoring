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

var telemetryCmd = &cobra.Command{
	Use:   "get-telemetry",
	Short: "Get telemetry for equipment",
	RunE:  getTelemetry,
}

func getTelemetry(cmd *cobra.Command, args []string) error {

	log.Debug().Msg("Getting Telemetry")

	config, err := getTelemetryConfig()
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
		var serialNumbers []string

		if config.DiscoverSerials {
			inventory, err := solaredgeClient.GetSiteInventory(siteID)
			if err != nil {
				return fmt.Errorf("unable to get inventory for site %s: %w", siteID, err)
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
			equipment, err := solaredgeClient.GetTelemetryForEquipment(siteID, serial, config.TimeUnit, config.StartTime, config.EndTime)
			if err != nil {
				return err
			}

			equipmentJSON, err := json.Marshal(equipment)
			if err != nil {
				return fmt.Errorf("unable to generate json for equipment serial %s at site %s: %w", serial, siteID, err)
			}

			filename := fmt.Sprintf("%s/%s_%s.json", config.OutputDir, siteID, serial)
			err = ioutil.WriteFile(filename, equipmentJSON, 0644)
			if err != nil {
				return fmt.Errorf("unable to write equipment file %s: %w", filename, err)
			}

			log.Debug().Interface("equipment", equipment).Msg("Got equipment telemetry")
		}

	}

	return nil
}

type telemetryConfig struct {
	StartTime time.Time
	EndTime   time.Time
	TimeUnit  string

	SiteIDs       []string
	SerialNumbers []string

	DiscoverSites   bool
	DiscoverSerials bool

	OutputDir string
}

func getTelemetryConfig() (*telemetryConfig, error) {
	config := telemetryConfig{}
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

	discoverSerials := viper.GetBool("all-equipment")
	serials := viper.GetStringSlice("serial-number")
	if discoverSerials {
		if len(serials) > 0 {
			errs = multierror.Append(errs, errors.New("cannot set all-equipment and specify serials"))
		} else {
			config.DiscoverSerials = true
		}
	} else {
		if len(serials) == 0 {
			errs = multierror.Append(errs, errors.New("must set all-equipment or specify at least one serial"))
		}
		config.SerialNumbers = serials
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

/*
func parseDate(timestamp string) (time.Time, error) {
	return time.Parse(api.DateFormat, timestamp)
}
*/

func parseTime(timestamp string) (time.Time, error) {
	return time.Parse(api.TimeFormat, timestamp)
}

func init() {
	RootCmd.AddCommand(telemetryCmd)

	telemetryCmd.Flags().StringP("start-time", "", "", "Specify the start time for telemetry - will default to 24 hours ago.")
	telemetryCmd.Flags().StringP("end-time", "", "", "Specify the start time for telemetry - will default to now.")
	telemetryCmd.Flags().BoolP("by-hour", "", false, "Specify hourly samples")
	telemetryCmd.Flags().BoolP("by-quarter-hour", "", false, "Specify 15-minute samples")
	telemetryCmd.Flags().BoolP("by-day", "", false, "Specify daily samples")

	telemetryCmd.Flags().StringSliceP("site-id", "", []string{}, "Specify site IDs; use multiple flags for multiple sites")
	telemetryCmd.Flags().BoolP("all-sites", "", false, "Discover available sites and use them all")
	telemetryCmd.Flags().StringSliceP("serial-number", "", []string{}, "Specify telemetry source serial numbers")
	telemetryCmd.Flags().BoolP("all-equipment", "", false, "Discover available equipment at each specified site")

	telemetryCmd.Flags().StringP("output-dir", "", "", "Specify where output files belong")
}
