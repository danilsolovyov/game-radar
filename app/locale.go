package app

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
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
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("GetUserDefaultLocaleName")

	buf := make([]uint16, 85)
	ret, _, _ := proc.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if ret == 0 {
		return "en-US"
	}
	return syscall.UTF16ToString(buf)
}
