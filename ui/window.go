package ui

import (
	"path/filepath"

	"altstash/cfg"

	libcoin "altstash/lib/taler/coin"
	libconfig "altstash/lib/config"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// Window is the main application window.
type Window struct {
	window       *adw.ApplicationWindow
	navView      *adw.NavigationView
	toastOverlay *adw.ToastOverlay

	carousel    *adw.Carousel
	bottomBar   *gtk.Box
	receiveBtn  *gtk.ToggleButton
	haveBtn     *gtk.ToggleButton
	sendBtn     *gtk.ToggleButton
	receivePage *ReceivePage
	sendPage    *SendPage
	updatingNav bool

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

	// Create NavigationView for the Have tab
	receiver.navView = adw.NewNavigationView()
	receiver.navView.SetHExpand(true)
	receiver.navView.SetVExpand(true)
	receiver.navView.Add(receiver.balanceList.page)

	// Create placeholder pages
	receiver.receivePage = newReceivePage()
	receiver.sendPage = newSendPage()

	// Create carousel with 3 pages
	receiver.carousel = adw.NewCarousel()
	receiver.carousel.SetAllowMouseDrag(true)
	receiver.carousel.SetVExpand(true)
	receiver.carousel.SetHExpand(true)
	receiver.carousel.Append(receiver.receivePage.widget)
	receiver.carousel.Append(receiver.navView)
	receiver.carousel.Append(receiver.sendPage.widget)

	// Create bottom bar with 3 toggle buttons
	receiver.receiveBtn = gtk.NewToggleButtonWithLabel("Receive")
	receiver.receiveBtn.SetHExpand(true)

	receiver.haveBtn = gtk.NewToggleButtonWithLabel("Have")
	receiver.haveBtn.SetHExpand(true)
	receiver.haveBtn.SetGroup(receiver.receiveBtn)

	receiver.sendBtn = gtk.NewToggleButtonWithLabel("Send")
	receiver.sendBtn.SetHExpand(true)
	receiver.sendBtn.SetGroup(receiver.receiveBtn)

	// Default: Have button active
	receiver.haveBtn.SetActive(true)

	receiver.bottomBar = gtk.NewBox(gtk.OrientationHorizontal, 0)
	receiver.bottomBar.Append(receiver.receiveBtn)
	receiver.bottomBar.Append(receiver.haveBtn)
	receiver.bottomBar.Append(receiver.sendBtn)

	// Outer vertical box: carousel + bottom bar
	outerBox := gtk.NewBox(gtk.OrientationVertical, 0)
	outerBox.Append(receiver.carousel)
	outerBox.Append(receiver.bottomBar)

	receiver.toastOverlay = adw.NewToastOverlay()
	receiver.toastOverlay.SetChild(outerBox)

	// Scroll carousel to page 1 (Have) initially
	receiver.carousel.ScrollTo(receiver.navView, false)

	// Wire carousel <-> bottom bar sync
	receiver.wireCarouselNav()

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

// wireCarouselNav connects the carousel page-changed signal and bottom bar
// toggle buttons to keep them in sync.
func (receiver *Window) wireCarouselNav() {
	buttons := []*gtk.ToggleButton{receiver.receiveBtn, receiver.haveBtn, receiver.sendBtn}

	receiver.carousel.ConnectPageChanged(func(index uint) {
		if index >= uint(len(buttons)) {
			return
		}
		receiver.updatingNav = true
		buttons[index].SetActive(true)
		receiver.updatingNav = false
	})

	wireToggle := func(btn *gtk.ToggleButton, target gtk.Widgetter) {
		btn.ConnectToggled(func() {
			if !btn.Active() {
				return
			}
			if receiver.updatingNav {
				return
			}
			receiver.carousel.ScrollTo(target, true)
		})
	}

	wireToggle(receiver.receiveBtn, receiver.receivePage.widget)
	wireToggle(receiver.haveBtn, receiver.navView)
	wireToggle(receiver.sendBtn, receiver.sendPage.widget)
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
