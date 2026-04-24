package models

import "fmt"

// Blip represents an audio device
//
//export Blip
type Blip struct {
	Angle     float64 `json:"angle"`
	Distance  float64 `json:"distance"`
	Intensity float64 `json:"intensity"`
}

func (b *Blip) IsNil() bool {
	return b.Angle == 0 && b.Distance == 0 && b.Intensity == 0
}

func (b *Blip) String() string {
	return fmt.Sprintf("%v %v %v", b.Angle, b.Distance, b.Intensity)
}
