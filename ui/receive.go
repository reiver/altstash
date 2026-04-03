package ui

import (
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// ReceivePage is a placeholder for the future Receive screen.
type ReceivePage struct {
	widget *adw.ToolbarView
}

// newReceivePage creates the Receive placeholder page.
func newReceivePage() *ReceivePage {
	var receiver ReceivePage

	icon := gtk.NewImageFromIconName("plus-large-square-outline-symbolic")
	icon.SetPixelSize(128)
	icon.SetHExpand(true)
	icon.SetVExpand(true)
	icon.SetHAlign(gtk.AlignCenter)
	icon.SetVAlign(gtk.AlignCenter)

	header := adw.NewHeaderBar()
	header.SetTitleWidget(adw.NewWindowTitle("Receive", ""))

	receiver.widget = adw.NewToolbarView()
	receiver.widget.AddTopBar(header)
	receiver.widget.SetContent(icon)
	receiver.widget.SetHExpand(true)
	receiver.widget.SetVExpand(true)

	return &receiver
}
