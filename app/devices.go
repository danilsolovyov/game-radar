package app

import (
	"game-radar/models"
)

// GetAudioDevices returns a list of audio devices that have at least 8 channels.
func (a *App) GetAudioDevices() ([]models.DeviceSpeakers, error) {
	devices, err := a.driver.DevicesSpeakers()
	if err != nil {
		return nil, err
	}
	var result []models.DeviceSpeakers
	for _, device := range devices {
		if device.Channels < 8 {
			continue
		}
		result = append(result, device)
	}
	return result, nil

}

// SetDevice sets the active audio device.
func (a *App) SetDevice(device models.DeviceSpeakers) error {
	a.stopRadar()
	err := a.config.setDevice(device)
	if err != nil {
		return err
	}
	a.startRadar()
	return nil
}

// GetSelectedDevice returns the active audio device.
func (a *App) GetSelectedDevice() (models.DeviceSpeakers, error) {
	return a.driver.GetDevice(a.config.radarDeviceSpeakersID())
}
