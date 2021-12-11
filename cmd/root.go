package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dreamlibrarian/solaredge-monitoring/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var solaredgeClient *client.Client

var RootCmd = &cobra.Command{
	Use:   "solaredge-monitoring",
	Short: "A toolkit for rendering data from Solaredge monitoring.",
	Long:  "A toolkit for rendering data from Solaredge monitoring.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			return fmt.Errorf("failed to bind flags for cmd %s: %w", cmd.Name(), err)
		}

		configFile := viper.GetString("config")

		viper.SetConfigName(configFile)
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.SetEnvPrefix("SOLAREDGE")
		viper.AutomaticEnv()
		viper.EnvKeyReplacer(strings.NewReplacer("-", "_"))

		err = viper.ReadInConfig()
		if err != nil {
			return fmt.Errorf("config parsing failed for cmd %s: %w", cmd.Name(), err)
		}

		verbose := viper.GetBool("verbose")
		if verbose {
			viper.Debug()
		}

		apiKey := viper.GetString("api-key")
		if apiKey == "" {
			return errors.New("api-key must be specified")
		}
		solaredgeClient = client.NewClient(apiKey)

		return nil
	},
}

func init() {

	RootCmd.PersistentFlags().StringP("api-key", "", "", "API Key")
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose Mode")
	RootCmd.PersistentFlags().StringP("config", "c", "solaredge.yml", "Config File")
}
