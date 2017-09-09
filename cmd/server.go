package cmd

import (
	"time"

	"github.com/donothingloop/hamgo/node"
	"github.com/donothingloop/hamgo/parameters"
	"github.com/donothingloop/hamgo/protocol"
	"github.com/donothingloop/hamgo/rest"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var test bool

func init() {
	serverCmd.Flags().BoolVar(&test, "test", false, "Enable testing")
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start the gossip server",
	Run:   executeServer,
}

func executeServer(cmd *cobra.Command, args []string) {
	// read config
	config = parameters.ReadConfig(configFile)

	sett := config.Node

	n, err := node.NewNode(sett)

	if err != nil {
		logrus.WithError(err).Warn("Failed to create node")
	}

	logrus.Info("Node started.")
	defer n.Close()

	// create a new rest server
	rs := rest.NewServer(config.REST)
	go rs.Init(n)

	if test {
		logrus.Debug("Waiting for 5 seconds")
		<-time.After(time.Second * 5)

		for {
			<-time.After(time.Second)
			n.SpreadMessage(&protocol.Message{
				Version: 1,
			})
		}
	}

	select {}
}
