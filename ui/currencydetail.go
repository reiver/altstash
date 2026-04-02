package ui

import (
	libcoin "altstash/lib/taler/coin"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// CurrencyDetailPage is Screen 2: per-exchange breakdown for a single currency.
type CurrencyDetailPage struct {
	page    *adw.NavigationPage
	listBox *gtk.ListBox

	balances []libcoin.ExchangeBalance
}

// newCurrencyDetailPage creates the currency detail page.
func newCurrencyDetailPage(currency string, balances []libcoin.ExchangeBalance) *CurrencyDetailPage {
	var receiver CurrencyDetailPage

	receiver.balances = balances

	receiver.listBox = gtk.NewListBox()
	receiver.listBox.SetSelectionMode(gtk.SelectionNone)
	receiver.listBox.AddCSSClass("boxed-list")

	for _, balance := range balances {
		row := adw.NewActionRow()
		row.SetTitle(balance.ExchangeBaseURL)

		icon := gtk.NewImageFromIconName("network-server-symbolic")
		row.AddPrefix(icon)

		amountLabel := gtk.NewLabel(balance.Total.FormatValue())
		row.AddSuffix(amountLabel)

		receiver.listBox.Append(row)
	}

	group := adw.NewPreferencesGroup()
	group.Add(receiver.listBox)

	if 0 == len(balances) {
		label := gtk.NewLabel("No exchanges for this currency")
		label.AddCSSClass("dim-label")
		label.SetMarginTop(24)
		label.SetMarginBottom(24)
		group.Add(label)
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

	receiver.page = adw.NewNavigationPage(toolbar, currency)

	return &receiver
}
