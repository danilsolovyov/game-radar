package models

import "image/color"

type Theme struct {
	Name                 string     `toml:"name" json:"name"`
	BackgroundColor      color.RGBA `toml:"background_color" json:"background_color"`
	RadarColor           color.RGBA `toml:"radar_color" json:"radar_color"`
	BorderOpacity        float64    `toml:"border_opacity" json:"border_opacity"`
	BorderWidth          int        `toml:"border_width" json:"border_width"`
	SectionBaseOpacity   float64    `toml:"section_base_opacity" json:"section_base_opacity"`
	SectionBrightOpacity float64    `toml:"section_bright_opacity" json:"section_bright_opacity"`
	SectionTimeout       int        `toml:"section_timeout" json:"section_timeout"`
	SectionCount         int        `toml:"section_count" json:"section_count"`
	RingCount            int        `toml:"ring_count" json:"ring_count"`
	ShowBlips            bool       `toml:"show_blips" json:"show_blips"`
	BlipOpacity          float64    `toml:"blip_opacity" json:"blip_opacity"`
	BlipTimeout          int        `toml:"blip_timeout" json:"blip_timeout"`
	BlipSize             int        `toml:"blip_size" json:"blip_size"`
	Size                 int        `toml:"size" json:"size"`
	PosX                 int        `toml:"pos_x" json:"pos_x"`
	PosY                 int        `toml:"pos_y" json:"pos_y"`
	IntensityMultiplier  float32    `toml:"intensity_multiplier" json:"intensity_multiplier"`
}

func DefaultTheme() Theme {
	return Theme{
		Name:                 "Default",
		BackgroundColor:      color.RGBA{0, 20, 10, 50},
		RadarColor:           color.RGBA{0, 255, 180, 255},
		BorderOpacity:        0.1,
		BorderWidth:          2,
		SectionBaseOpacity:   0,
		SectionBrightOpacity: 1,
		SectionTimeout:       250,
		RingCount:            3,
		SectionCount:         25,
		ShowBlips:            true,
		BlipOpacity:          0.5,
		BlipTimeout:          250,
		BlipSize:             3,
		Size:                 320,
		PosX:                 30,
		PosY:                 30,
		IntensityMultiplier:  2,
	}
}

// OrangeTheme returns the second built-in theme from the default set.
func OrangeTheme() Theme {
	return Theme{
		Name:                 "Orange",
		BackgroundColor:      color.RGBA{20, 10, 0, 100},
		RadarColor:           color.RGBA{255, 123, 0, 255},
		BorderOpacity:        0.1,
		BorderWidth:          2,
		SectionBaseOpacity:   0,
		SectionBrightOpacity: 0.9,
		SectionTimeout:       500,
		RingCount:            3,
		SectionCount:         25,
		ShowBlips:            true,
		BlipOpacity:          0.5,
		BlipTimeout:          500,
		BlipSize:             3,
		Size:                 320,
		PosX:                 500,
		PosY:                 30,
		IntensityMultiplier:  1,
	}
}

// BuiltinThemes returns application built-in themes.
func BuiltinThemes() map[string]Theme {
	defaultTheme := DefaultTheme()
	orangeTheme := OrangeTheme()
	return map[string]Theme{
		defaultTheme.Name: defaultTheme,
		orangeTheme.Name:  orangeTheme,
	}
}
