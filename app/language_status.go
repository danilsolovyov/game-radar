package app

import wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

const languageChangedEventName = "language-changed"

func (a *App) emitLanguageChanged(language string) {
	if a.ctx == nil {
		return
	}
	wailsRuntime.EventsEmit(a.ctx, languageChangedEventName, language)
}
