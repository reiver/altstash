package ui

import (
	"os"

	"altstash/cfg"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
)

// Run starts the application main loop.
func Run() int {
	app := adw.NewApplication(cfg.AppID, gio.ApplicationFlagsNone)
	app.ConnectActivate(func() {
		onActivate(app)
	})
	return app.Run(os.Args)
}

func onActivate(app *adw.Application) {
	win := app.ActiveWindow()
	if nil == win {
		w := newWindow(app)
		w.window.Present()
		return
	}
	win.Present()
}
