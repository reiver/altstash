package ui

import (
	_ "embed"
	"os"
	"path/filepath"

	"altstash/cfg"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

//go:embed icons/wallet-symbolic.svg
var walletIconSVG []byte

//go:embed icons/link.reiver.altstash.svg
var appIconSVG []byte

// installEmbeddedIcons extracts embedded icons to the XDG cache directory
// and registers the path with GTK's icon theme so that
// gtk.NewImageFromIconName can find them.
func installEmbeddedIcons() error {
	iconDir := filepath.Join(cfg.IconsDir(), "hicolor", "scalable", "actions")

	err := os.MkdirAll(iconDir, 0755)
	if nil != err {
		return err
	}

	err = os.WriteFile(filepath.Join(iconDir, "wallet-symbolic.svg"), walletIconSVG, 0644)
	if nil != err {
		return err
	}

	// extract app icon
	appsDir := filepath.Join(cfg.IconsDir(), "hicolor", "scalable", "apps")

	err = os.MkdirAll(appsDir, 0755)
	if nil != err {
		return err
	}

	err = os.WriteFile(filepath.Join(appsDir, "link.reiver.altstash.svg"), appIconSVG, 0644)
	if nil != err {
		return err
	}

	// register icon theme search path
	iconTheme := gtk.IconThemeGetForDisplay(gdk.DisplayGetDefault())
	iconTheme.AddSearchPath(cfg.IconsDir())

	return nil
}
