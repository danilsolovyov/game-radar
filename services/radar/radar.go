package radar

import (
	"context"
	"game-radar/models"
	"game-radar/services/audiodriver"
	"log/slog"
	"math"
)

const (
	fl  = iota // front left | 330
	fr         // front right 30
	fc         // front center | 0
	lwe        // low frequency effects | 0 not in use
	bl         // back left (rear left) | 210
	br         // back right (rear right) | 150
	sl         // side left | 270
	sr         // side right | 90

	channelSize = 6
	historySize = 6
)

type AudioDriver interface {
	LoopbackDevice(ctx context.Context, device models.DeviceSpeakers, reader audiodriver.SpeakersPeaksWriter) error
}

type SpeakerChannel struct {
	Name      string
	Angle     float64
	PeakValue float32
	IsUnused  bool
}

func GetSpeakerChannelsMap71(peaks []float32) []SpeakerChannel {
	if len(peaks) < 8 {
		return nil
	}

	return []SpeakerChannel{
		{
			Name:      "front_left",
			Angle:     330,
			PeakValue: peaks[fl],
		},
		{
			Name:      "front_right",
			Angle:     30,
			PeakValue: peaks[fr],
		},
		{
			Name:      "front_center",
			Angle:     0,
			PeakValue: peaks[fc],
		},
		{
			Name:      "low_frequency_effects",
			Angle:     0,
			PeakValue: peaks[lwe],
			IsUnused:  true,
		},
		{
			Name:      "back_left",
			Angle:     210,
			PeakValue: peaks[bl],
		},
		{
			Name:      "back_right",
			Angle:     150,
			PeakValue: peaks[br],
		},
		{
			Name:      "side_left",
			Angle:     270,
			PeakValue: peaks[sl],
		},
		{
			Name:      "side_right",
			Angle:     90,
			PeakValue: peaks[sr],
		},
	}
}

type Radar struct {
	audioDriver     audiodriver.AudioDriver
	peaksChan       chan []float32
	blipsChan       chan models.Blip
	history         [][]float32
	intensityFilter float32
	amplifier       float64
	logger          *slog.Logger
}

func NewRadarService(audioDriver audiodriver.AudioDriver, intensityFilter float32, amplifier float64, logger *slog.Logger) *Radar {
	r := &Radar{
		audioDriver:     audioDriver,
		blipsChan:       make(chan models.Blip, channelSize),
		history:         make([][]float32, historySize),
		intensityFilter: intensityFilter,
		amplifier:       amplifier,
		logger:          logger,
	}
	return r
}

func (r *Radar) Start(ctx context.Context, device models.DeviceSpeakers) {
	r.peaksChan = make(chan []float32, channelSize)
	go func() {
		_ = r.audioDriver.LoopbackDeviceSpeakers(ctx, device, r)
	}()

	// Process data coming from the audio driver.
	go r.processAudioData(ctx, device)
}

func (r *Radar) Write(peaks []float32) {
	r.peaksChan <- peaks
}

func (r *Radar) processAudioData(ctx context.Context, device models.DeviceSpeakers) {
	var iterator int
	for peaks := range r.peaksChan {
		select {
		case <-ctx.Done():
			return
		default:
			r.history[iterator] = peaks

			// fill history
			if iterator < len(r.history)-1 {
				iterator++
				continue
			}

			weightedHistory := r.calculateWeightedHistory(device.Channels)

			blip := r.calculateBlip(GetSpeakerChannelsMap71(weightedHistory))
			if !blip.IsNil() {
				r.blipsChan <- blip
			}

			r.history = make([][]float32, historySize)
			iterator = 0
		}
	}
}

func (r *Radar) calculateWeightedHistory(channelsCount uint32) []float32 {
	weightedPeaks := make([]float32, channelsCount)

	historySize := len(r.history)
	if historySize == 0 {
		return weightedPeaks
	}

	for _, peaks := range r.history {
		for i := range channelsCount {
			if int(i) >= len(peaks) {
				break
			}
			weightedPeaks[i] += peaks[i]
		}
	}

	for i := range weightedPeaks {
		weightedPeaks[i] /= float32(historySize)
	}

	return weightedPeaks
}

func (r *Radar) calculateBlip(channels []SpeakerChannel) models.Blip {
	var sinSum, cosSum float64
	panChannels := 0

	for _, channel := range channels {
		if channel.IsUnused {
			continue
		}

		panChannels++
		if channel.PeakValue < r.intensityFilter {
			continue
		}

		angle := (channel.Angle * math.Pi) / 180
		sinSum += math.Sin(angle) * float64(channel.PeakValue)
		cosSum += math.Cos(angle) * float64(channel.PeakValue)
	}

	if panChannels == 0 {
		return models.Blip{}
	}

	count := float64(panChannels)
	sinSum /= count
	cosSum /= count

	blip := models.Blip{
		Angle:     math.Atan2(sinSum, cosSum) * 180 / math.Pi,
		Intensity: math.Sqrt(cosSum*cosSum+sinSum*sinSum) * r.amplifier,
	}

	blip.Angle = math.Mod(blip.Angle+360, 360)
	blip.Angle = math.Round(blip.Angle)
	if blip.Angle == 360 {
		blip.Angle = 0
	}

	if blip.Intensity > 0.001 {
		r.logger.Debug("channels", "channels", channels, "blip", blip)
	}

	return blip
}

func (r *Radar) GetBlips() <-chan models.Blip {
	return r.blipsChan
}
