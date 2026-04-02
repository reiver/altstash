package ui

import (
	"context"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// SettingsPage is the settings/preferences window.
type SettingsPage struct {
	window *adw.PreferencesWindow

	// OnDataDirectoryChanged is called when the user selects a new data directory.
	OnDataDirectoryChanged func(newPath string)
}

// newSettingsPage creates the settings window.
func newSettingsPage(parent *gtk.Window, currentDataDir string) *SettingsPage {
	var receiver SettingsPage

	dirRow := adw.NewActionRow()
	dirRow.SetTitle("Data Directory")
	dirRow.SetSubtitle(currentDataDir)

	dialog := gtk.NewFileDialog()
	dialog.SetTitle("Choose Data Directory")

	changeBtn := gtk.NewButtonWithLabel("Change")
	changeBtn.SetVAlign(gtk.AlignCenter)
	changeBtn.ConnectClicked(func() {
		dialog.SelectFolder(context.Background(), parent, func(result gio.AsyncResulter) {
			file, err := dialog.SelectFolderFinish(result)
			if nil != err {
				return
			}

			newPath := file.Path()
			if "" == newPath {
				return
			}

			dirRow.SetSubtitle(newPath)

			if nil != receiver.OnDataDirectoryChanged {
				receiver.OnDataDirectoryChanged(newPath)
			}
		})
	})
	dirRow.AddSuffix(changeBtn)

	group := adw.NewPreferencesGroup()
	group.SetTitle("Storage")
	group.Add(dirRow)

	page := adw.NewPreferencesPage()
	page.Add(group)

	receiver.window = adw.NewPreferencesWindow()
	receiver.window.SetTitle("Settings")
	receiver.window.SetTransientFor(parent)
	receiver.window.Add(page)

	return &receiver
}
