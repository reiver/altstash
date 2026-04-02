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

// LoadConfigDir reads the config from configDir/config.json and returns a [Config].
// The config data is expected to be JSON.
//
//
// If the config data does not contain a "data_directory" field, or the field
// has an empty value, then the value is set to `defaultDataDir`.
//
// If the file doesn't exist, returns a Config with defaultDataDir.
//
// See also: [LoadFromBytes].
func LoadConfigDir(configDir string, defaultDataDir string) (Config, error) {
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

	return LoadFromBytes(data, defaultDataDir)
}

// Load reads the config from a []byte and returns a [Config].
// The config data is expected to be JSON.
//
// If the config data does not contain a "data_directory" field, or the field
// has an empty value, then the value is set to `defaultDataDir`.
//
// See also: [LoadConfigDir].
func LoadFromBytes(bytes []byte, defaultDataDir string) (Config, error) {
	var config Config

	err := json.Unmarshal(bytes, &config)
	if nil != err {
		return Config{}, erorr.Wrap(err, "could not parse config JSON")
	}

	if "" == config.DataDirectory {
		config.DataDirectory = defaultDataDir
	}

	return config, nil
}
