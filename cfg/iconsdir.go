package cfg

import (
	"path/filepath"

	"codeberg.org/reiver/go-env"
)

// IconsDir returns $XDG_CACHE_HOME/altstash/icons (typically ~/.cache/altstash/icons).
func IconsDir() string {
	cacheHome := env.GetElse[string]("XDG_CACHE_HOME", filepath.Join(userHomeDir, ".cache"))
	return filepath.Join(cacheHome, "altstash", "icons")
}
