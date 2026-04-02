package cfg

import (
	"path/filepath"

	"codeberg.org/reiver/go-env"
)

// ConfigDir returns $XDG_CONFIG_HOME/altstash/ (typically ~/.config/altstash/).
func ConfigDir() string {
	configHome := env.GetElse[string]("XDG_CONFIG_HOME", filepath.Join(userHomeDir, ".config"))
	return filepath.Join(configHome, "altstash")
}
