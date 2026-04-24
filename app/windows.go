package app

import (
	"errors"
	"fmt"
	"syscall"

	"github.com/lxn/win"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	// popupStyle is WS_POPUP represented as int32 (high bit set).
	popupStyle int32 = -2147483648

	// Extended style masks used when switching between overlay and normal modes.
	overlayAddExStyleMask = win.WS_EX_LAYERED | win.WS_EX_TRANSPARENT |
		win.WS_EX_TOOLWINDOW
	overlayClearExStyleMask = win.WS_EX_APPWINDOW
	normalAddExStyleMask    = win.WS_EX_APPWINDOW
	normalClearExStyleMask  = win.WS_EX_TRANSPARENT | win.WS_EX_TOOLWINDOW

	// Base style and SetWindowPos flags used for frame updates and z-order changes.
	overlayClearStyleMask = win.WS_OVERLAPPEDWINDOW | win.WS_CAPTION | win.WS_THICKFRAME |
		win.WS_BORDER | win.WS_DLGFRAME | win.WS_VSCROLL | win.WS_HSCROLL
	frameChangeFlags = win.SWP_NOMOVE | win.SWP_NOSIZE | win.SWP_NOZORDER |
		win.SWP_NOACTIVATE | win.SWP_FRAMECHANGED
	topmostFlags = win.SWP_NOMOVE | win.SWP_NOSIZE
)

// SetRadarOverlay switches the app window into radar overlay mode.
// It removes window decorations, applies transparent/click-through extended styles,
// places the window above others, and synchronizes overlay status.
func (a *App) SetRadarOverlay() {
	if err := a.setWindowStyle(true); err != nil {
		a.logger.ErrorContext(a.ctx, "Error setting frameless overlay style", "error", err)
		return
	}

	err := a.setExtendedWindowStyle(
		overlayAddExStyleMask,
		overlayClearExStyleMask,
	)
	if err != nil {
		a.logger.ErrorContext(a.ctx, "Error setting window long", "error", err)
		return
	}

	wailsRuntime.WindowSetAlwaysOnTop(a.ctx, true)
	wailsRuntime.WindowSetSize(a.ctx, a.config.getTheme().Size, a.config.getTheme().Size)
	wailsRuntime.WindowSetPosition(a.ctx, a.config.getTheme().PosX, a.config.getTheme().PosY)

	a.setTopmost(true)
	a.UpdateOverlayStatus(true)
}

// SetNormalWindow restores the app window to the regular desktop mode.
// It reapplies standard window decorations, disables overlay extended styles,
// restores size/position behavior, and synchronizes overlay status.
func (a *App) SetNormalWindow() {
	if err := a.setWindowStyle(false); err != nil {
		a.logger.ErrorContext(a.ctx, "Error restoring normal window style", "error", err)
		return
	}

	err := a.setExtendedWindowStyle(
		normalAddExStyleMask,
		normalClearExStyleMask,
	)
	if err != nil {
		a.logger.ErrorContext(a.ctx, "Error setting window long", "error", err)
		return
	}
	wailsRuntime.WindowSetAlwaysOnTop(a.ctx, false)
	wailsRuntime.WindowSetSize(a.ctx, DefaultWindowWidth, DefaultWindowHeight)
	wailsRuntime.WindowCenter(a.ctx)
	a.setTopmost(false)
	a.UpdateOverlayStatus(false)
}

// getHwnd resolves the native window handle by the current application title.
func (a *App) getHwnd() (win.HWND, error) {
	lpWindowName, err := syscall.UTF16PtrFromString(a.appName)
	if err != nil {
		return 0, fmt.Errorf("error converting app name to UTF16Ptr: %w", err)
	}
	hwnd := win.FindWindow(nil, lpWindowName)
	if hwnd == 0 {
		return 0, errors.New("window not found")
	}
	return hwnd, nil
}

// setExtendedWindowStyle updates GWL_EXSTYLE by adding and clearing provided flags.
func (a *App) setExtendedWindowStyle(addFlags int32, clearFlags int32) error {
	hwnd, err := a.getHwnd()
	if err != nil {
		return err
	}

	currentExStyle := win.GetWindowLong(hwnd, win.GWL_EXSTYLE)
	nextExStyle := (currentExStyle | addFlags) &^ clearFlags
	return applyWindowStyle(hwnd, win.GWL_EXSTYLE, nextExStyle)
}

// setWindowStyle updates GWL_STYLE to either frameless popup or decorated window mode.
func (a *App) setWindowStyle(frameless bool) error {
	hwnd, err := a.getHwnd()
	if err != nil {
		return err
	}

	currentStyle := win.GetWindowLong(hwnd, win.GWL_STYLE)
	newStyle := currentStyle

	if frameless {
		// Remove all system chrome (caption, borders, and scrollbars).
		newStyle &^= overlayClearStyleMask
		newStyle |= popupStyle
	} else {
		// Restore default overlapped window style and clear popup flag.
		newStyle &^= popupStyle
		newStyle |= win.WS_OVERLAPPEDWINDOW
	}

	applyErr := applyWindowStyle(hwnd, win.GWL_STYLE, newStyle)
	if applyErr != nil {
		return applyErr
	}

	if frameless {
		// Force redraw path to avoid stale title bar artifacts after frame change.
		win.ShowWindow(hwnd, win.SW_SHOW)
	}
	return nil
}

// applyWindowStyle writes a style value (GWL_STYLE or GWL_EXSTYLE) and commits it.
// The SWP_FRAMECHANGED call is required so Win10/Win11 immediately reapplies chrome.
func applyWindowStyle(hwnd win.HWND, index int32, value int32) error {
	win.SetWindowLong(hwnd, index, value)
	// Force immediate frame recalculation; otherwise UI updates may be delayed.
	if !win.SetWindowPos(hwnd, 0, 0, 0, 0, 0, frameChangeFlags) {
		return fmt.Errorf("failed to apply window frame changes: %w", syscall.GetLastError())
	}
	return nil
}

// setTopmost toggles topmost z-order state for the current window.
func (a *App) setTopmost(enabled bool) {
	hwnd, err := a.getHwnd()
	if err != nil {
		return
	}

	insertAfter := win.HWND_NOTOPMOST
	if enabled {
		insertAfter = win.HWND_TOPMOST
	}
	win.SetWindowPos(hwnd, insertAfter, 0, 0, 0, 0, topmostFlags)
}
