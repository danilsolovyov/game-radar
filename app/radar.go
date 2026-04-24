package app

import (
	"context"
	"errors"
	"fmt"
	"game-radar/models"
	"runtime/debug"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	startupRadarInitTimeout = 4 * time.Second
	getDeviceTimeout        = 3 * time.Second
)

func (a *App) startRadarAsync() {
	done := make(chan struct{})
	go func() {
		defer close(done)
		a.startRadar()
	}()

	go func() {
		select {
		case <-done:
			a.logger.InfoContext(a.ctx, "startup radar init finished")
		case <-time.After(startupRadarInitTimeout):
			a.logger.WarnContext(a.ctx, "startup radar init is slow; app continues without blocking UI", "timeout", startupRadarInitTimeout.String())
			<-done
			a.logger.InfoContext(a.ctx, "startup radar init finished after timeout")
		}
	}()
}

func (a *App) startRadar() {
	defer func() {
		if r := recover(); r != nil {
			a.logger.ErrorContext(a.ctx, "startRadar panic recovered", "panic", r, "stack", string(debug.Stack()))
		}
	}()

	// Reset previous session: keep a single blips reader and cancel the prior radar pipeline.
	a.logger.InfoContext(a.ctx, "startRadar begin")
	a.stopRadar()

	device, err := a.getDevice()
	if err != nil {
		a.logger.ErrorContext(a.ctx, "startRadar aborted: failed to resolve audio device", "error", err)
		return
	}
	if device == (models.DeviceSpeakers{}) {
		a.logger.ErrorContext(a.ctx, "startRadar aborted: default audio device not found")
		return
	}

	radarCtx, cancel := context.WithCancel(a.ctx)
	a.radarCancelFn = cancel
	a.logger.InfoContext(
		a.ctx,
		"startRadar device selected",
		"device_id", device.ID,
		"device_name", device.Name,
		"channels", device.Channels,
		"is_default", device.IsDefault,
	)
	a.radar.Start(radarCtx, device)

	go a.forwardRadarData(radarCtx)
	a.logger.InfoContext(a.ctx, "startRadar complete")
}

func (a *App) getDevice() (models.DeviceSpeakers, error) {
	if id := a.config.radarDeviceSpeakersID(); id != "" {
		a.logger.InfoContext(a.ctx, "getDevice: trying configured device id", "device_id", id)
		device, err := a.getConfiguredDeviceWithTimeout(id, getDeviceTimeout)
		if err != nil {
			if errors.Is(err, errGetDeviceTimeout) {
				a.logger.ErrorContext(a.ctx, "getDevice: configured device lookup timed out", "device_id", id, "timeout", getDeviceTimeout.String())
			} else {
				a.logger.ErrorContext(a.ctx, "getDevice: configured device lookup failed", "device_id", id, "error", err)
			}
			return models.DeviceSpeakers{}, err
		}
		return device, nil
	}

	a.logger.InfoContext(a.ctx, "getDevice: no configured device id, enumerating active speakers")
	devices, err := a.getDefaultDeviceFromEnumerationWithTimeout(getDeviceTimeout)
	if err != nil {
		if errors.Is(err, errGetDeviceTimeout) {
			a.logger.ErrorContext(a.ctx, "getDevice: speaker enumeration timed out", "timeout", getDeviceTimeout.String())
		} else {
			a.logger.ErrorContext(a.ctx, "getDevice: speaker enumeration failed", "error", err)
		}
		return models.DeviceSpeakers{}, err
	}
	return devices, nil
}

var errGetDeviceTimeout = errors.New("audio device call timeout")

func (a *App) getConfiguredDeviceWithTimeout(deviceID string, timeout time.Duration) (models.DeviceSpeakers, error) {
	type result struct {
		device models.DeviceSpeakers
		err    error
	}
	ch := make(chan result, 1)
	go func() {
		device, err := a.driver.GetDevice(deviceID)
		ch <- result{device: device, err: err}
	}()

	select {
	case out := <-ch:
		return out.device, out.err
	case <-time.After(timeout):
		return models.DeviceSpeakers{}, fmt.Errorf("%w: GetDevice(%s)", errGetDeviceTimeout, deviceID)
	}
}

func (a *App) getDefaultDeviceFromEnumerationWithTimeout(timeout time.Duration) (models.DeviceSpeakers, error) {
	type result struct {
		device models.DeviceSpeakers
		err    error
	}
	ch := make(chan result, 1)
	go func() {
		devices, err := a.driver.DevicesSpeakers()
		if err != nil {
			ch <- result{err: err}
			return
		}
		a.logger.InfoContext(a.ctx, "getDevice: active speakers enumerated", "count", len(devices))

		for _, d := range devices {
			if d.IsDefault {
				a.logger.InfoContext(a.ctx, "getDevice: default speaker found", "device_id", d.ID, "device_name", d.Name)
				ch <- result{device: d}
				return
			}
		}
		a.logger.WarnContext(a.ctx, "getDevice: default speaker not found in active list", "count", len(devices))
		ch <- result{device: models.DeviceSpeakers{}}
	}()

	select {
	case out := <-ch:
		return out.device, out.err
	case <-time.After(timeout):
		return models.DeviceSpeakers{}, fmt.Errorf("%w: DevicesSpeakers()", errGetDeviceTimeout)
	}
}

// forwardRadarData forwards blips to frontend until radarCtx is canceled (single channel consumer).
func (a *App) forwardRadarData(radarCtx context.Context) {
	ch := a.radar.GetBlips()
	for {
		select {
		case <-radarCtx.Done():
			return
		case data := <-ch:
			if radarCtx.Err() != nil {
				return
			}
			if a.ctx != nil {
				wailsRuntime.EventsEmit(a.ctx, "radar-data", data)
			}
		}
	}
}

func (a *App) stopRadar() {
	if a.radarCancelFn != nil {
		a.radarCancelFn()
		a.radarCancelFn = nil
	}
}

// StopRadar stops capture and radar event bridge (for shutdown and device switch).
func (a *App) StopRadar() {
	a.stopRadar()
}

// GetBlipType returns a Blip type for Wails to generate TypeScript model.
func (a *App) GetBlipType() models.Blip {
	return models.Blip{}
}

func (a *App) GetTheme() models.Theme {
	return a.config.getTheme()
}
