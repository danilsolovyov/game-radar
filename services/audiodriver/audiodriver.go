package audiodriver

import (
	"context"
	"game-radar/models"
)

type AudioDriver interface {
	DevicesSpeakers() ([]models.DeviceSpeakers, error)
	LoopbackDeviceSpeakers(ctx context.Context, device models.DeviceSpeakers, writer SpeakersPeaksWriter) error
	GetDevice(id string) (models.DeviceSpeakers, error)
}
