package audiodriver

import (
	"context"
	"errors"
	"fmt"
	"game-radar/models"
	"log/slog"
	"runtime"
	"time"
	"unsafe"

	"github.com/danilsolovyov/go-wca/pkg/wca"
	"github.com/go-ole/go-ole"
)

type SpeakersPeaksWriter interface {
	Write([]float32)
}

type WCA struct {
	logger *slog.Logger
}

func initCOM() (func(), error) {
	if err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED); err == nil {
		return ole.CoUninitialize, nil
	} else {
		// Fallback: some environments already run in STA.
		if errSta := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); errSta == nil {
			return ole.CoUninitialize, nil
		} else {
			return nil, fmt.Errorf("coInitializeEx failed (MTA: %w, STA: %w)", err, errSta)
		}
	}
}

func NewWCA(logger *slog.Logger) *WCA {
	if logger == nil {
		logger = slog.Default()
	}
	return &WCA{logger: logger}
}

func (w *WCA) LoopbackDeviceSpeakers(ctx context.Context, device models.DeviceSpeakers, writer SpeakersPeaksWriter) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	uninit, err := initCOM()
	if err != nil {
		return err
	}
	defer uninit()

	var de *wca.IMMDeviceEnumerator
	if err := wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &de); err != nil {
		return fmt.Errorf("CoCreateInstance failed: %w", err)
	}
	defer de.Release()

	// TODO: implement go-wca enumerator that retrieves a device by ID.
	// Enumerate all devices and find the required one by ID.
	var devCollection *wca.IMMDeviceCollection
	if err := de.EnumAudioEndpoints(wca.ERender, wca.DEVICE_STATE_ACTIVE, &devCollection); err != nil {
		return fmt.Errorf("EnumAudioEndpoints failed: %w", err)
	}
	defer devCollection.Release()

	var count uint32
	if err := devCollection.GetCount(&count); err != nil {
		return fmt.Errorf("GetCount failed: %w", err)
	}

	var mmd *wca.IMMDevice
	for i := range count {
		var currentDevice *wca.IMMDevice
		if err := devCollection.Item(i, &currentDevice); err != nil {
			continue
		}

		var deviceID string
		if err := currentDevice.GetId(&deviceID); err != nil {
			currentDevice.Release()
			continue
		}

		if deviceID == device.ID {
			mmd = currentDevice
			break
		}
		// currentDevice.Release()
	}

	if mmd == nil {
		return fmt.Errorf("device with ID %s not found", device.ID)
	}
	defer mmd.Release()

	var ps *wca.IPropertyStore
	if err := mmd.OpenPropertyStore(wca.STGM_READ, &ps); err != nil {
		return fmt.Errorf("openPropertyStore failed: %w", err)
	}
	defer ps.Release()

	var rac *wca.IAudioClient
	if err := mmd.Activate(wca.IID_IAudioClient, wca.CLSCTX_ALL, nil, &rac); err != nil {
		return fmt.Errorf("activate IAudioClient failed: %w", err)
	}
	defer rac.Release()

	var wfx *wca.WAVEFORMATEX
	if err := rac.GetMixFormat(&wfx); err != nil {
		return fmt.Errorf("getMixFormat failed: %w", err)
	}
	defer ole.CoTaskMemFree(uintptr(unsafe.Pointer(wfx)))

	var ami *wca.IAudioMeterInformation
	if err := mmd.Activate(wca.IID_IAudioMeterInformation, wca.CLSCTX_ALL, nil, &ami); err != nil {
		return fmt.Errorf("activate IAudioMeterInformation failed: %w", err)
	}
	defer ami.Release()

	// Use device period from configuration.

	if err := rac.Initialize(wca.AUDCLNT_SHAREMODE_SHARED,
		wca.AUDCLNT_STREAMFLAGS_EVENTCALLBACK|wca.AUDCLNT_STREAMFLAGS_LOOPBACK,
		wca.REFERENCE_TIME(device.DefaultPeriod),
		0,
		wfx,
		nil,
	); err != nil {
		return fmt.Errorf("initialize failed: %w", err)
	}

	audioReadyEvent := wca.CreateEventExA(0, 0, 0, wca.EVENT_MODIFY_STATE|wca.SYNCHRONIZE)
	defer wca.CloseHandle(audioReadyEvent)

	if err := rac.SetEventHandle(audioReadyEvent); err != nil {
		return fmt.Errorf("setEventHandle failed: %w", err)
	}

	if err := rac.Start(); err != nil {
		return fmt.Errorf("start failed: %w", err)
	}
	defer func() {
		if stopErr := rac.Stop(); stopErr != nil {
			w.logger.Error("stop failed", "error", stopErr)
		}
	}()

	peaks := make([]float32, device.Channels)

	// Small delay for initialization.
	initSleep := 300 * time.Millisecond
	time.Sleep(initSleep)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := watchEvent(ctx, audioReadyEvent); err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}
				return fmt.Errorf("watchEvent failed: %w", err)
			}

			if err := ami.GetChannelsPeakValues(device.Channels, peaks); err != nil {
				return fmt.Errorf("getChannelsPeakValues failed: %w", err)
			}
			writer.Write(peaks)

			// Delay between polling cycles.
			if device.Latency > 0 {
				time.Sleep(device.Latency)
			} else {
				time.Sleep(50 * time.Millisecond) // default delay
			}
		}
	}
}

func (w *WCA) DevicesSpeakers() ([]models.DeviceSpeakers, error) {
	w.logger.Info("wca.DevicesSpeakers: begin")
	var devices []models.DeviceSpeakers

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	w.logger.Info("wca.DevicesSpeakers: os thread locked")

	uninit, err := initCOM()
	if err != nil {
		w.logger.Error("wca.DevicesSpeakers: initCOM failed", "error", err)
		return nil, err
	}
	w.logger.Info("wca.DevicesSpeakers: COM initialized")
	defer uninit()

	var de *wca.IMMDeviceEnumerator
	w.logger.Info("wca.DevicesSpeakers: creating MMDeviceEnumerator")
	if err := wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &de); err != nil {
		w.logger.Error("wca.DevicesSpeakers: CoCreateInstance failed", "error", err)
		return nil, fmt.Errorf("coCreateInstance failed: %w", err)
	}
	w.logger.Info("wca.DevicesSpeakers: MMDeviceEnumerator created")
	defer de.Release()

	var devCollection *wca.IMMDeviceCollection
	w.logger.Info("wca.DevicesSpeakers: calling EnumAudioEndpoints")
	if err := de.EnumAudioEndpoints(wca.ERender, wca.DEVICE_STATE_ACTIVE, &devCollection); err != nil {
		w.logger.Error("wca.DevicesSpeakers: EnumAudioEndpoints failed", "error", err)
		return nil, fmt.Errorf("enumAudioEndpoints failed: %w", err)
	}
	w.logger.Info("wca.DevicesSpeakers: EnumAudioEndpoints ok")
	defer devCollection.Release()

	var count uint32
	w.logger.Info("wca.DevicesSpeakers: reading device count")
	if err := devCollection.GetCount(&count); err != nil {
		w.logger.Error("wca.DevicesSpeakers: GetCount failed", "error", err)
		return nil, fmt.Errorf("getCount failed: %w", err)
	}
	w.logger.Info("wca.DevicesSpeakers: device count read", "count", count)

	defaultDeviceID, err := w.getDefaultDeviceID(de)
	if err != nil {
		w.logger.Error("wca.DevicesSpeakers: getDefaultDeviceID failed", "error", err)
		return nil, fmt.Errorf("getDefaultDeviceID failed: %w", err)
	}
	w.logger.Info("wca.DevicesSpeakers: default device id resolved", "device_id", defaultDeviceID)

	for i := range count {
		var dev *wca.IMMDevice
		if err = devCollection.Item(i, &dev); err != nil {
			w.logger.Warn("failed to get device", "index", i, "error", err)
			continue
		}

		var device models.DeviceSpeakers
		device, err = w.getDeviceInfo(dev)
		if err != nil {
			w.logger.Warn("failed to get device info", "index", i, "error", err)
			dev.Release()
			continue
		}

		device.IsDefault = device.ID == defaultDeviceID

		devices = append(devices, device)
		dev.Release()
	}

	w.logger.Info("wca.DevicesSpeakers: completed", "devices", len(devices))
	return devices, nil
}

func (w *WCA) GetDevice(id string) (models.DeviceSpeakers, error) {
	w.logger.Info("wca.GetDevice: begin", "device_id", id)
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	w.logger.Info("wca.GetDevice: os thread locked", "device_id", id)

	uninit, err := initCOM()
	if err != nil {
		w.logger.Error("wca.GetDevice: initCOM failed", "device_id", id, "error", err)
		return models.DeviceSpeakers{}, err
	}
	w.logger.Info("wca.GetDevice: COM initialized", "device_id", id)
	defer uninit()

	var de *wca.IMMDeviceEnumerator
	w.logger.Info("wca.GetDevice: creating MMDeviceEnumerator", "device_id", id)
	if err := wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &de); err != nil {
		w.logger.Error("wca.GetDevice: CoCreateInstance failed", "device_id", id, "error", err)
		return models.DeviceSpeakers{}, fmt.Errorf("coCreateInstance failed: %w", err)
	}
	w.logger.Info("wca.GetDevice: MMDeviceEnumerator created", "device_id", id)
	defer de.Release()

	var dev *wca.IMMDevice
	w.logger.Info("wca.GetDevice: calling IMMDeviceEnumerator.GetDevice", "device_id", id)
	if err := de.GetDevice(id, &dev); err != nil {
		w.logger.Error("wca.GetDevice: IMMDeviceEnumerator.GetDevice failed", "device_id", id, "error", err)
		return models.DeviceSpeakers{}, fmt.Errorf("getDevice failed: %w", err)
	}
	w.logger.Info("wca.GetDevice: IMMDeviceEnumerator.GetDevice ok", "device_id", id)
	defer dev.Release()

	device, err := w.getDeviceInfo(dev)
	if err != nil {
		w.logger.Error("wca.GetDevice: getDeviceInfo failed", "device_id", id, "error", err)
		return models.DeviceSpeakers{}, err
	}
	w.logger.Info("wca.GetDevice: completed", "device_id", id, "device_name", device.Name, "channels", device.Channels)
	return device, nil
}

func (w *WCA) getDeviceInfo(dev *wca.IMMDevice) (models.DeviceSpeakers, error) {
	var device models.DeviceSpeakers

	// Read device ID.
	var pwszID string
	if err := dev.GetId(&pwszID); err != nil {
		return device, fmt.Errorf("getId failed: %w", err)
	}
	device.ID = pwszID

	// Read device properties.
	var props *wca.IPropertyStore
	if err := dev.OpenPropertyStore(wca.STGM_READ, &props); err != nil {
		return device, fmt.Errorf("openPropertyStore failed: %w", err)
	}
	defer props.Release()

	// Read device friendly name.
	var pv wca.PROPVARIANT
	if err := props.GetValue(&wca.PKEY_Device_FriendlyName, &pv); err != nil {
		return device, fmt.Errorf("getValue PKEY_Device_FriendlyName failed: %w", err)
	}

	device.Name = pv.String()

	// Activate device to read audio format.
	var audioClient *wca.IAudioClient
	if err := dev.Activate(wca.IID_IAudioClient, wca.CLSCTX_ALL, nil, &audioClient); err != nil {
		return device, fmt.Errorf("activate IAudioClient failed: %w", err)
	}
	defer audioClient.Release()

	// Read wave format.
	var mixFormat *wca.WAVEFORMATEX
	if err := audioClient.GetMixFormat(&mixFormat); err != nil {
		return device, fmt.Errorf("getMixFormat failed: %w", err)
	}
	defer ole.CoTaskMemFree(uintptr(unsafe.Pointer(mixFormat)))

	// Fill format details.
	device.FormatPCM = mixFormat.WBitsPerSample
	device.Rate = mixFormat.NSamplesPerSec
	// device.Channels = uint32(mixFormat.NChannels)

	// Read device period details.
	var defaultPeriod, minimumPeriod wca.REFERENCE_TIME
	if err := audioClient.GetDevicePeriod(&defaultPeriod, &minimumPeriod); err != nil {
		return device, fmt.Errorf("getDevicePeriod failed: %w", err)
	}
	device.DefaultPeriod = int64(defaultPeriod)
	device.MinimumPeriod = int64(minimumPeriod)
	device.Latency = time.Duration(defaultPeriod)

	// Read metering details.
	var meter *wca.IAudioMeterInformation
	if err := dev.Activate(wca.IID_IAudioMeterInformation, wca.CLSCTX_ALL, nil, &meter); err == nil {
		defer meter.Release()
		var meterChannels uint32
		if err = meter.GetMeteringChannelCount(&meterChannels); err == nil {
			device.Channels = meterChannels
		}
	}

	return device, nil
}

func (w *WCA) getDefaultDeviceID(de *wca.IMMDeviceEnumerator) (string, error) {
	var defaultDev *wca.IMMDevice
	if err := de.GetDefaultAudioEndpoint(wca.ERender, wca.EConsole, &defaultDev); err != nil {
		return "", fmt.Errorf("getDefaultAudioEndpoint failed: %w", err)
	}
	defer defaultDev.Release()

	var pwszID string
	if err := defaultDev.GetId(&pwszID); err != nil {
		return "", fmt.Errorf("getId for default device failed: %w", err)
	}

	return pwszID, nil
}

func watchEvent(ctx context.Context, event uintptr) error {
	result := make(chan error, 1)

	go func() {
		dw := wca.WaitForSingleObject(event, wca.INFINITE)
		if dw != 0 {
			result <- fmt.Errorf("WaitForSingleObject failed with code %d", dw)
		} else {
			result <- nil
		}
	}()

	select {
	case err := <-result:
		return err
	case <-ctx.Done():
		return context.Canceled
	}
}
