package libconfig

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"codeberg.org/reiver/go-erorr"
)

// Config holds the user's altstash configuration.
type Config struct {
	DataDirectory string `json:"data_directory"`
}

// Load reads the config from configDir/config.json.
// If the file doesn't exist, returns a Config with defaultDataDir.
// The caller provides the paths so that lib/ does not import cfg/.
func Load(configDir string, defaultDataDir string) (Config, error) {
	configPath := filepath.Join(configDir, "config.json")

	data, err := os.ReadFile(configPath)
	if nil != err {
		if errors.Is(err, os.ErrNotExist) {
			return Config{
				DataDirectory: defaultDataDir,
			}, nil
		}
		return Config{}, erorr.Wrap(err, "could not read config file")
	}

	var config Config

	err = json.Unmarshal(data, &config)
	if nil != err {
		return Config{}, erorr.Wrap(err, "could not parse config JSON")
	}

	if "" == config.DataDirectory {
		config.DataDirectory = defaultDataDir
	}

	return config, nil
}
