package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/dreamlibrarian/solaredge-monitoring/action"
	"github.com/dreamlibrarian/solaredge-monitoring/api"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var telemetryCmd = &cobra.Command{
	Use:   "get-telemetry",
	Short: "Get telemetry for equipment",
	RunE: func(cmd *cobra.Command, args []string) error {

		config, err := getTelemetryConfig()
		if err != nil {
			return err
		}

		action := action.NewTelemetryAction(apiKey)

		fsMap, err := action.Do(config)
		if err != nil {
			return err
		}

		for path, contents := range fsMap {
			err = ioutil.WriteFile(path, contents, 0644)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

/*
func parseDate(timestamp string) (time.Time, error) {
	return time.Parse(api.DateFormat, timestamp)
}
*/

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

func getTelemetryConfig() (*action.TelemetryActionConfig, error) {
	config := action.TelemetryActionConfig{}
	var err, errs error

	startTime := viper.GetString("start-time")
	if startTime == "" {
		config.StartTime = time.Now().Add(-24 * time.Hour)
	} else {
		if config.StartTime, err = api.ParseTime(startTime); err != nil {
			errs = multierror.Append(errs, fmt.Errorf("start time could not be parsed: %w", err))
		}
	}

	endTime := viper.GetString("end-time")
	if endTime == "" {
		config.EndTime = time.Now()
	} else {
		if config.EndTime, err = api.ParseTime(endTime); err != nil {
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

	config.DiscoverSites = viper.GetBool("all-sites")
	config.SiteIDs = viper.GetStringSlice("site-id")

	config.DiscoverSerials = viper.GetBool("all-equipment")
	config.SerialNumbers = viper.GetStringSlice("serial-number")

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
