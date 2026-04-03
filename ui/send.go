package ui

import (
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// SendPage is a placeholder for the future Send screen.
type SendPage struct {
	widget *adw.ToolbarView
}

// newSendPage creates the Send placeholder page.
func newSendPage() *SendPage {
	var receiver SendPage

	icon := gtk.NewImageFromIconName("minus-square-outline-symbolic")
	icon.SetPixelSize(128)
	icon.SetHExpand(true)
	icon.SetVExpand(true)
	icon.SetHAlign(gtk.AlignCenter)
	icon.SetVAlign(gtk.AlignCenter)

	header := adw.NewHeaderBar()
	header.SetTitleWidget(adw.NewWindowTitle("Send", ""))

	receiver.widget = adw.NewToolbarView()
	receiver.widget.AddTopBar(header)
	receiver.widget.SetContent(icon)
	receiver.widget.SetHExpand(true)
	receiver.widget.SetVExpand(true)

	return &receiver
}
