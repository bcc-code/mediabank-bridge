package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/bcc-code/mediabank-bridge/log"
	"github.com/bcc-code/mediabank-bridge/vantage"
)

var defaultConfigPaths = []string{
	"~/.config/mediabank-bridge/config.json",
	".config.json",
	"config.json",
}

// Config structure for the app read from a JSON file
type Config struct {
	Vantage vantage.ClientSettings
}

// Validate if all parts of the config are ok
func (c Config) Validate() error {
	return c.Vantage.Validate()
}

// MustReadConfigFile either from the passed path or from one of the default locations
func MustReadConfigFile(filePath string) Config {
	if filePath != "" {
		if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
			log.L.Fatal().
				Str("filePath", filePath).
				Msg("Specified config file does not exist. Did *not* attempt to look in default locations")
		}
	} else {
		for _, fp := range defaultConfigPaths {
			log.L.Debug().Str("filePath", fp).Msg("Looking for config file")
			if _, err := os.Stat(fp); errors.Is(err, os.ErrNotExist) {
				continue
			}
			filePath = fp
			break
		}
	}

	if filePath == "" {
		log.L.Fatal().
			Msg("Unable to find config in any paths")
	}

	fp, err := os.Open(filePath)
	if err != nil {
		log.L.Fatal().Err(err)
	}

	byteValue, err := ioutil.ReadAll(fp)
	if err != nil {
		log.L.Fatal().Err(err)
	}

	cfg := Config{}
	err = json.Unmarshal(byteValue, &cfg)
	if err != nil {
		log.L.Fatal().Err(err)
	}

	err = cfg.Validate()
	if err != nil {
		log.L.Fatal().Err(err)
	}

	return cfg
}
