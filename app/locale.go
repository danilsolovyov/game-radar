package app

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func normalizeAppLanguage(language string) string {
	normalized := strings.ToLower(strings.TrimSpace(language))
	switch {
	case strings.HasPrefix(normalized, "ru"):
		return "ru"
	case strings.HasPrefix(normalized, "en"):
		return "en"
	default:
		return "en"
	}
}

func (a *App) SetLanguage(language string) error {
	normalized := normalizeAppLanguage(language)
	if err := a.config.SetLanguage(normalized); err != nil {
		return fmt.Errorf("app.SetLanguage: failed to save config: %w", err)
	}
	a.refreshTrayTexts()
	a.emitLanguageChanged(normalized)
	return nil
}

func (a *App) GetLanguage() string {
	return a.config.getLanguage()
}

func GetSystemLocaleWindows() string {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Control Panel\International`, registry.QUERY_VALUE)
	if err != nil {
		return "en-US"
	}
	defer k.Close()
	locale, _, err := k.GetStringValue("LocaleName")
	if err != nil || locale == "" {
		return "en-US"
	}
	return locale
}
