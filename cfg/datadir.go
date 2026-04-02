package cfg

import (
	"path/filepath"

	"codeberg.org/reiver/go-env"
)

// DefaultDataDir returns $XDG_DATA_HOME/altstash/ (typically ~/.local/share/altstash/).
// This is where coin files are stored by default. The user can override via config file.
func DefaultDataDir() string {
	dataHome := env.GetElse[string]("XDG_DATA_HOME", filepath.Join(userHomeDir, ".local", "share"))
	return filepath.Join(dataHome, "altstash")
}
