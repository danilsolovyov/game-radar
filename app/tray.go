package app

import (
	"strings"

	"github.com/getlantern/systray"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type trayTexts struct {
	appTitle      string
	appTooltip    string
	quitTitle     string
	quitTooltip   string
	toggleOn      string
	toggleOnHint  string
	toggleOff     string
	toggleOffHint string
}

var trayI18n = map[string]trayTexts{
	"en": {
		appTitle:      "Game Radar",
		appTooltip:    "Game Radar",
		quitTitle:     "Quit",
		quitTooltip:   "Exit application",
		toggleOn:      "Open Overlay",
		toggleOnHint:  "Switch to overlay mode",
		toggleOff:     "Close Overlay",
		toggleOffHint: "Return to window mode",
	},
	"ru": {
		appTitle:      "Game Radar",
		appTooltip:    "Игровой звуковой радар",
		quitTitle:     "Выход",
		quitTooltip:   "Завершить приложение",
		toggleOn:      "Открыть оверлей",
		toggleOnHint:  "Переключить в режим оверлея",
		toggleOff:     "Закрыть оверлей",
		toggleOffHint: "Вернуться в оконный режим",
	},
}

func normalizeLanguageCode(language string) string {
	code := strings.ToLower(strings.TrimSpace(language))
	if strings.HasPrefix(code, "ru") {
		return "ru"
	}
	return "en"
}

func (a *App) trayTextsForCurrentLanguage() trayTexts {
	if text, ok := trayI18n[normalizeLanguageCode(a.config.getLanguage())]; ok {
		return text
	}
	return trayI18n["en"]
}

// InitTray configures system tray items and starts click handler.
func (a *App) InitTray(icon []byte) {
	a.mu.Lock()
	a.trayIcon = icon
	a.mu.Unlock()

	systray.SetIcon(icon)
	systray.SetTooltip(a.trayTitle)
	systray.SetTitle(a.trayTitle)

	a.mu.Lock()
	a.trayToggleItem = systray.AddMenuItem("", "")
	systray.AddSeparator()
	a.trayQuitItem = systray.AddMenuItem("", "")
	stopCh := a.trayStopCh
	a.mu.Unlock()

	a.refreshTrayTexts()

	go func() {
		for {
			select {
			case <-a.trayToggleItem.ClickedCh:
				a.ToggleWindowModeFromTray()
			case <-a.trayQuitItem.ClickedCh:
				a.QuitFromTray()
				systray.Quit()
				return
			case <-stopCh:
				return
			}
		}
	}()
}

// ShutdownTray stops the tray click handler.
func (a *App) ShutdownTray() {
	a.mu.Lock()
	defer a.mu.Unlock()
	select {
	case <-a.trayStopCh:
	default:
		close(a.trayStopCh)
	}
}

func (a *App) refreshTrayTexts() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.refreshTrayTextsLocked()
}

func (a *App) refreshTrayTextsLocked() {
	if a.trayToggleItem == nil || a.trayQuitItem == nil {
		return
	}
	text := a.trayTextsForCurrentLanguage()
	systray.SetTitle(text.appTitle)
	systray.SetTooltip(text.appTooltip)
	if a.overlayMode {
		a.trayToggleItem.SetTitle(text.toggleOff)
		a.trayToggleItem.SetTooltip(text.toggleOffHint)
	} else {
		a.trayToggleItem.SetTitle(text.toggleOn)
		a.trayToggleItem.SetTooltip(text.toggleOnHint)
	}
	a.trayQuitItem.SetTitle(text.quitTitle)
	a.trayQuitItem.SetTooltip(text.quitTooltip)
}

// ToggleWindowModeFromTray switches window mode between overlay and normal mode.
func (a *App) ToggleWindowModeFromTray() {
	if a.ctx == nil {
		return
	}
	if a.IsOverlayMode() {
		a.SetNormalWindow()
		return
	}
	a.SetRadarOverlay()
}

// QuitFromTray exits the application from a tray command.
func (a *App) QuitFromTray() {
	if a.ctx == nil {
		return
	}
	wailsRuntime.Quit(a.ctx)
}
