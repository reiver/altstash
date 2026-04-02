package ui

import (
	"altstash/cfg"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// showAboutDialog presents the application about dialog.
func showAboutDialog(parent *gtk.Window) {
	about := adw.NewAboutWindow()
	about.SetTransientFor(parent)
	about.SetApplicationName(cfg.Name)
	about.SetApplicationIcon(cfg.AppID)
	about.SetVersion(cfg.Version)
	about.SetDeveloperName(cfg.AuthorName)
	about.SetWebsite(cfg.AuthorWebSite)
	about.SetLicenseType(gtk.LicenseMITX11)
	about.SetCopyright(cfg.CopyRightMessage)
	about.SetComments(cfg.TagLine)
	about.Present()
}
