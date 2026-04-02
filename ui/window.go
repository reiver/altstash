package ui

import (
	"path/filepath"

	"altstash/cfg"

	libcoin "altstash/lib/taler/coin"
	libconfig "altstash/lib/config"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
)

// Window is the main application window.
type Window struct {
	window       *adw.ApplicationWindow
	navView      *adw.NavigationView
	toastOverlay *adw.ToastOverlay

	balanceList *BalanceListPage
	config      libconfig.Config

	coins []libcoin.Coin
}

// newWindow creates and configures the main application window.
func newWindow(app *adw.Application) *Window {
	var receiver Window

	// install embedded icons
	iconsErr := installEmbeddedIcons()

	// load config and coins
	var configErr error
	receiver.config, configErr = libconfig.LoadConfigDir(cfg.ConfigDir(), cfg.DefaultDataDir())
	var coinsErr error
	talerCoinsDir := filepath.Join(receiver.config.DataDirectory, cfg.TalerCoinsDir)
	receiver.coins, coinsErr = libcoin.LoadFromDirectory(talerCoinsDir)

	balances, balancesErr := libcoin.BalanceByCurrency(receiver.coins)

	receiver.balanceList = newBalanceListPage(balances)
	receiver.wireBalanceListCallbacks()

	receiver.navView = adw.NewNavigationView()
	receiver.navView.Add(receiver.balanceList.page)

	receiver.toastOverlay = adw.NewToastOverlay()
	receiver.toastOverlay.SetChild(receiver.navView)

	receiver.window = adw.NewApplicationWindow(&app.Application)
	receiver.window.SetTitle(cfg.Name)
	receiver.window.SetDefaultSize(360, 648)
	receiver.window.SetContent(receiver.toastOverlay)

	// show error toasts if something failed
	if nil != iconsErr {
		receiver.toastOverlay.AddToast(adw.NewToast("Could not install icons: " + iconsErr.Error()))
	}
	if nil != configErr {
		receiver.toastOverlay.AddToast(adw.NewToast("Could not load config: " + configErr.Error()))
	}
	if nil != coinsErr {
		receiver.toastOverlay.AddToast(adw.NewToast("Could not load coins: " + coinsErr.Error()))
	}
	if nil != balancesErr {
		receiver.toastOverlay.AddToast(adw.NewToast("Could not compute balances: " + balancesErr.Error()))
	}

	return &receiver
}

// wireBalanceListCallbacks sets all callbacks on the current balanceList.
// Called from newWindow and refreshBalances to avoid duplication.
func (receiver *Window) wireBalanceListCallbacks() {
	receiver.balanceList.OnCurrencyActivated = func(currency string) {
		exchangeBalances, err := libcoin.BalanceByCurrencyAndExchange(receiver.coins, currency)
		if nil != err {
			receiver.toastOverlay.AddToast(adw.NewToast("Could not compute balances: " + err.Error()))
			return
		}
		detailPage := newCurrencyDetailPage(currency, exchangeBalances)
		receiver.navView.Push(detailPage.page)
	}
	receiver.balanceList.OnAbout = func() {
		showAboutDialog(&receiver.window.Window)
	}
	receiver.balanceList.OnSettings = func() {
		settingsPage := newSettingsPage(&receiver.window.Window, receiver.config.DataDirectory)
		settingsPage.OnDataDirectoryChanged = func(newPath string) {
			receiver.config.DataDirectory = newPath
			err := libconfig.Save(cfg.ConfigDir(), receiver.config)
			if nil != err {
				receiver.toastOverlay.AddToast(adw.NewToast("Could not save config: " + err.Error()))
				return
			}
			receiver.refreshBalances()
		}
		settingsPage.window.Present()
	}
}

// refreshBalances reloads coins and rebuilds the balance list page.
func (receiver *Window) refreshBalances() {
	talerCoinsDir := filepath.Join(receiver.config.DataDirectory, cfg.TalerCoinsDir)

	var coinsErr error
	receiver.coins, coinsErr = libcoin.LoadFromDirectory(talerCoinsDir)
	if nil != coinsErr {
		receiver.toastOverlay.AddToast(adw.NewToast("Could not load coins: " + coinsErr.Error()))
	}

	balances, balancesErr := libcoin.BalanceByCurrency(receiver.coins)
	if nil != balancesErr {
		receiver.toastOverlay.AddToast(adw.NewToast("Could not compute balances: " + balancesErr.Error()))
	}

	receiver.navView.Remove(receiver.balanceList.page)
	receiver.balanceList = newBalanceListPage(balances)
	receiver.wireBalanceListCallbacks()
	receiver.navView.Add(receiver.balanceList.page)
}
