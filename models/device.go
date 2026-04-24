package models

import (
	"time"
)

// DeviceSpeakers represents an audio device
//
//export DeviceSpeakers
type DeviceSpeakers struct {
	ID            string        `json:"id" toml:"id"`                         // id
	Name          string        `json:"name" toml:"name"`                     // audio device name
	FormatPCM     uint16        `json:"format_pcm" toml:"format_pcm"`         // wbits per sample
	Rate          uint32        `json:"rate" toml:"rate"`                     // nsamples per second
	Channels      uint32        `json:"channels" toml:"channels"`             // channels
	DefaultPeriod int64         `json:"default_period" toml:"default_period"` // default period in ns
	MinimumPeriod int64         `json:"minimum_period" toml:"minimum_period"` // minimum period in ns
	Latency       time.Duration `json:"latency" toml:"latency"`               // calculated latency: DefaultPeriod * 100
	IsDefault     bool          `json:"is_default" toml:"is_default"`         // whether the device is the default device
}

func (d *DeviceSpeakers) CalculateLatency() {
	d.Latency = time.Duration(d.DefaultPeriod * 100)
}
