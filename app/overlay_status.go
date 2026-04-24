package app

import wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

const overlayStatusEventName = "overlay-status"

func (a *App) startOverlayStatusBroadcaster() {
	go func() {
		for status := range a.overlayStatus {
			if a.ctx == nil {
				continue
			}
			wailsRuntime.EventsEmit(a.ctx, overlayStatusEventName, status)
		}
	}()
}

func (a *App) emitOverlayStatus(status bool) {
	select {
	case a.overlayStatus <- status:
	default:
		select {
		case <-a.overlayStatus:
		default:
		}
		a.overlayStatus <- status
	}
}

// UpdateOverlayStatus synchronizes overlay status.
func (a *App) UpdateOverlayStatus(status bool) {
	a.mu.Lock()
	a.overlayMode = status
	a.refreshTrayTextsLocked()
	a.mu.Unlock()
	a.emitOverlayStatus(status)
}

// IsOverlayMode returns current window mode.
func (a *App) IsOverlayMode() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.overlayMode
}
