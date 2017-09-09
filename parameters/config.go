package parameters

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
)

// Config stores the configuration for the servers.
type Config struct {
	REST    RESTSettings `json:"rest"`
	Node    Settings     `json:"node"`
	Station Station      `json:"station"`
}

// ReadConfig reads the configuration JSON file.
func ReadConfig(file string) *Config {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to read config file %s", file)
		return nil
	}

	cfg := Config{}

	// unmarshal the config
	if err := json.Unmarshal(dat, &cfg); err != nil {
		logrus.WithError(err).Fatal("Failed to parse config")
	}

	return &cfg
}
