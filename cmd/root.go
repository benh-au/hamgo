package cmd

import (
	"github.com/donothingloop/hamgo/parameters"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hamgo",
	Short: "Resource and user discovery for HAMNET based on a gossip protocol",
	Long:  "hamgo - Resource and user discovery for HAMNET based on a gossip protocol",
}

var debug bool
var config *parameters.Config
var configFile string

// Execute the root cmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("execution failed")
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug output")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "hamgo.json", "JSON config file")
}

func initConfig() {
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debugging output enabled")
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}
