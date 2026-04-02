package ui

import (
	libcoin "altstash/lib/taler/coin"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"altstash/cfg"
)

// BalanceListPage is Screen 1: combined per-currency balances.
type BalanceListPage struct {
	page    *adw.NavigationPage
	listBox *gtk.ListBox

	balances []libcoin.CurrencyBalance

	// OnCurrencyActivated is called when the user taps a currency row.
	OnCurrencyActivated func(currency string)
}

// newBalanceListPage creates the balance list page.
func newBalanceListPage(balances []libcoin.CurrencyBalance) *BalanceListPage {
	var receiver BalanceListPage

	receiver.balances = balances

	receiver.listBox = gtk.NewListBox()
	receiver.listBox.SetSelectionMode(gtk.SelectionNone)
	receiver.listBox.AddCSSClass("boxed-list")

	receiver.listBox.ConnectRowActivated(func(row *gtk.ListBoxRow) {
		if nil == receiver.OnCurrencyActivated {
			return
		}

		index := row.Index()
		if index < 0 || index >= len(receiver.balances) {
			return
		}

		receiver.OnCurrencyActivated(receiver.balances[index].Currency)
	})

	group := adw.NewPreferencesGroup()
	group.Add(receiver.listBox)

	if 0 == len(balances) {
		// empty state
		label := gtk.NewLabel("No coins yet")
		label.AddCSSClass("dim-label")
		label.SetMarginTop(24)
		label.SetMarginBottom(24)
		group.Add(label)
	} else {
		for _, balance := range balances {
			row := adw.NewActionRow()
			row.SetTitle(balance.Currency)
			row.SetActivatable(true)

			// We do a bit of a trick to make this icon available.
			//
			// The original SVG source-code for it is at:
			// altstash/ui/icons/wallet-symbolic.svg
			//
			// But, we have to put it on the file system to get GNOME/GTK to use it.
			//
			// The code here just uses it.
			icon := gtk.NewImageFromIconName("wallet-symbolic")
			row.AddPrefix(icon)

			amountLabel := gtk.NewLabel(balance.Total.FormatValue())
			row.AddSuffix(amountLabel)

			arrow := gtk.NewImageFromIconName("go-next-symbolic")
			row.AddSuffix(arrow)

			receiver.listBox.Append(row)
		}
	}

	contentBox := gtk.NewBox(gtk.OrientationVertical, 12)
	contentBox.SetMarginTop(12)
	contentBox.SetMarginBottom(12)
	contentBox.SetMarginStart(12)
	contentBox.SetMarginEnd(12)
	contentBox.Append(group)

	scrolled := gtk.NewScrolledWindow()
	scrolled.SetPolicy(gtk.PolicyNever, gtk.PolicyAutomatic)
	scrolled.SetVExpand(true)
	scrolled.SetChild(contentBox)

	clamp := adw.NewClamp()
	clamp.SetMaximumSize(600)
	clamp.SetChild(scrolled)

	header := adw.NewHeaderBar()

	toolbar := adw.NewToolbarView()
	toolbar.AddTopBar(header)
	toolbar.SetContent(clamp)

	receiver.page = adw.NewNavigationPage(toolbar, cfg.Name)

	return &receiver
}
