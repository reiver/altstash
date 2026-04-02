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

	// OnSettings is called when the user taps the "Settings" menu item.
	OnSettings func()

	// OnAbout is called when the user taps the "About" menu item.
	OnAbout func()
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

	// menu popover
	settingsRow := adw.NewActionRow()
	settingsRow.SetTitle("Settings")
	settingsRow.SetActivatable(true)
	settingsRow.AddPrefix(gtk.NewImageFromIconName("preferences-system-symbolic"))

	aboutRow := adw.NewActionRow()
	aboutRow.SetTitle("About " + cfg.Name)
	aboutRow.SetActivatable(true)
	aboutRow.AddPrefix(gtk.NewImageFromIconName("help-about-symbolic"))

	menuList := gtk.NewListBox()
	menuList.SetSelectionMode(gtk.SelectionNone)
	menuList.AddCSSClass("boxed-list")
	menuList.SetMarginTop(6)
	menuList.SetMarginBottom(6)
	menuList.SetMarginStart(6)
	menuList.SetMarginEnd(6)
	menuList.Append(settingsRow)
	menuList.Append(aboutRow)

	popover := gtk.NewPopover()
	popover.SetChild(menuList)

	menuList.ConnectRowActivated(func(row *gtk.ListBoxRow) {
		switch row.Index() {
		case 0:
			popover.Popdown()
			if nil != receiver.OnSettings {
				receiver.OnSettings()
			}
		case 1:
			popover.Popdown()
			if nil != receiver.OnAbout {
				receiver.OnAbout()
			}
		}
	})

	menuBtn := gtk.NewMenuButton()
	menuBtn.SetIconName("open-menu-symbolic")
	menuBtn.SetTooltipText("Menu")
	menuBtn.SetPopover(popover)

	header := adw.NewHeaderBar()
	header.PackEnd(menuBtn)

	toolbar := adw.NewToolbarView()
	toolbar.AddTopBar(header)
	toolbar.SetContent(clamp)

	receiver.page = adw.NewNavigationPage(toolbar, cfg.Name)

	return &receiver
}
