package app

import (
	"errors"
	"fmt"
	"game-radar/models"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
)

const (
	configFile          = "config.toml"
	DefaultWindowWidth  = 1400
	DefaultWindowHeight = 1100
	DefaultWindowPosX   = 0
	DefaultWindowPosY   = 0
	DefaultLogsDir      = "logs"
	DefaultAppLogFile   = "app.log"
	DefaultWailsLogFile = "wails.log"

	// DefaultIntensityFilter is the minimum channel peak (0..1) below which a channel is ignored in the direction vector.
	DefaultIntensityFilter float32 = 0.1
	// DefaultAmplifier scales blip intensity after analysis (before frontend theme processing).
	DefaultAmplifier float64 = 1.0
)

type Config struct {
	mu sync.RWMutex `toml:"-"`

	Language string `toml:"language"`
	Radar    struct {
		ThemeName        string  `toml:"theme_name"`
		DeviceSpeakersID string  `toml:"device_speakers_id"`
		IntensityFilter  float32 `toml:"intensity_filter"` // low-channel threshold; 0 in file -> DefaultIntensityFilter
		Amplifier        float64 `toml:"amplifier"`        // analysis gain; 0 in file -> DefaultAmplifier
	} `toml:"radar"`
	Logs struct {
		Dir     string `toml:"dir"`
		Backend struct {
			File  string `toml:"file"`
			Level string `toml:"level"`
		} `toml:"backend"`
		Wails struct {
			File            string `toml:"file"`
			Level           string `toml:"level"`
			LevelProduction string `toml:"level_production"`
		} `toml:"wails"`
		Rotation struct {
			MaxSizeMB  int  `toml:"max_size_mb"`
			MaxBackups int  `toml:"max_backups"`
			MaxAgeDays int  `toml:"max_age_days"`
			Compress   bool `toml:"compress"`
		} `toml:"rotation"`
	} `toml:"logs"`
	Themes map[string]models.Theme `toml:"themes"`
}

// saveConfig persists config to disk (thread-safe).
func (c *Config) saveConfig() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.persistUnlocked()
}

// persistUnlocked writes TOML; call only while c.mu is held (Lock).
func (c *Config) persistUnlocked() error {
	if c.Themes == nil {
		c.Themes = make(map[string]models.Theme)
	}

	file, err := os.OpenFile(configFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	if err = toml.NewEncoder(file).Encode(c); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}

// ensureRadarAnalysisDefaults fills radar analysis defaults if TOML fields are not set (zero after decode).
func (c *Config) ensureRadarAnalysisDefaults() {
	if c.Radar.IntensityFilter == 0 {
		c.Radar.IntensityFilter = DefaultIntensityFilter
	}
	if c.Radar.Amplifier == 0 {
		c.Radar.Amplifier = DefaultAmplifier
	}
}

// ensureLogsDefaults applies safe logging defaults.
func (c *Config) ensureLogsDefaults() {
	if strings.TrimSpace(c.Logs.Dir) == "" {
		c.Logs.Dir = DefaultLogsDir
	}
	if strings.TrimSpace(c.Logs.Backend.File) == "" {
		c.Logs.Backend.File = DefaultAppLogFile
	}
	if strings.TrimSpace(c.Logs.Backend.Level) == "" {
		c.Logs.Backend.Level = "info"
	}
	if strings.TrimSpace(c.Logs.Wails.File) == "" {
		c.Logs.Wails.File = DefaultWailsLogFile
	}
	if strings.TrimSpace(c.Logs.Wails.Level) == "" {
		c.Logs.Wails.Level = "info"
	}
	if strings.TrimSpace(c.Logs.Wails.LevelProduction) == "" {
		c.Logs.Wails.LevelProduction = "error"
	}
	if c.Logs.Rotation.MaxSizeMB <= 0 {
		c.Logs.Rotation.MaxSizeMB = 10
	}
	if c.Logs.Rotation.MaxBackups <= 0 {
		c.Logs.Rotation.MaxBackups = 5
	}
	if c.Logs.Rotation.MaxAgeDays <= 0 {
		c.Logs.Rotation.MaxAgeDays = 14
	}
	// Compress rotated archives by default when rotation section is effectively unset.
	if c.Logs.Rotation.MaxSizeMB == 10 &&
		c.Logs.Rotation.MaxBackups == 5 &&
		c.Logs.Rotation.MaxAgeDays == 14 &&
		!c.Logs.Rotation.Compress {
		c.Logs.Rotation.Compress = true
	}
}

// ensureBuiltinThemes guarantees built-in themes are present in config.
// User themes are preserved and never overwritten.
func (c *Config) ensureBuiltinThemes() {
	if c.Themes == nil {
		c.Themes = make(map[string]models.Theme)
	}

	for name, theme := range models.BuiltinThemes() {
		if _, exists := c.Themes[name]; !exists {
			c.Themes[name] = theme
		}
	}
}

func (c *Config) setDevice(device models.DeviceSpeakers) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	device.CalculateLatency()
	c.Radar.DeviceSpeakersID = device.ID
	return c.persistUnlocked()
}

// SetLanguage persists language code to config.
func (c *Config) SetLanguage(language string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Language = normalizeAppLanguage(language)
	return c.persistUnlocked()
}

func (c *Config) getTheme() models.Theme {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Themes == nil {
		c.Themes = make(map[string]models.Theme)
	}

	if theme, ok := c.Themes[c.Radar.ThemeName]; ok {
		return theme
	}

	defaultTheme := models.DefaultTheme()
	if theme, ok := c.Themes[defaultTheme.Name]; ok {
		c.Radar.ThemeName = defaultTheme.Name
		return theme
	}

	if len(c.Themes) > 0 {
		names := make([]string, 0, len(c.Themes))
		for name := range c.Themes {
			names = append(names, name)
		}
		sort.Strings(names)
		c.Radar.ThemeName = names[0]
		return c.Themes[names[0]]
	}

	c.Themes[defaultTheme.Name] = defaultTheme
	c.Radar.ThemeName = defaultTheme.Name
	return defaultTheme
}

func (c *Config) setTheme(theme models.Theme) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Themes == nil {
		c.Themes = make(map[string]models.Theme)
	}

	themeName := strings.TrimSpace(theme.Name)
	if themeName == "" {
		return errors.New("theme name cannot be empty")
	}

	theme.Name = themeName
	c.Themes[themeName] = theme
	c.Radar.ThemeName = themeName
	return c.persistUnlocked()
}

func (c *Config) getCurrentThemeName() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Radar.ThemeName
}

func (c *Config) getThemes() []models.Theme {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.Themes == nil {
		return []models.Theme{}
	}

	names := make([]string, 0, len(c.Themes))
	for name := range c.Themes {
		names = append(names, name)
	}
	sort.Strings(names)

	themes := make([]models.Theme, 0, len(names))
	for _, name := range names {
		themes = append(themes, c.Themes[name])
	}

	return themes
}

func (c *Config) deleteTheme(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Themes == nil {
		c.Themes = make(map[string]models.Theme)
	}

	themeName := strings.TrimSpace(name)
	if themeName == "" {
		return errors.New("theme name cannot be empty")
	}

	if _, ok := c.Themes[themeName]; !ok {
		return fmt.Errorf("theme %q not found", themeName)
	}

	if len(c.Themes) <= 1 {
		return errors.New("cannot delete the last theme")
	}

	delete(c.Themes, themeName)
	if c.Radar.ThemeName == themeName {
		defaultTheme := models.DefaultTheme()
		if _, ok := c.Themes[defaultTheme.Name]; ok {
			c.Radar.ThemeName = defaultTheme.Name
		} else {
			names := make([]string, 0, len(c.Themes))
			for theme := range c.Themes {
				names = append(names, theme)
			}
			sort.Strings(names)
			c.Radar.ThemeName = names[0]
		}
	}

	return c.persistUnlocked()
}

// radarDeviceSpeakersID returns selected device ID (thread-safe read).
func (c *Config) radarDeviceSpeakersID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Radar.DeviceSpeakersID
}

func (c *Config) getLanguage() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return normalizeAppLanguage(c.Language)
}

func LoadConfig() (*Config, error) {
	config := &Config{}
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		defaultTheme := models.DefaultTheme()
		config.Radar.ThemeName = defaultTheme.Name
		config.Themes = models.BuiltinThemes()
		config.ensureRadarAnalysisDefaults()
		config.ensureLogsDefaults()
		err = config.saveConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to save config: %w", err)
		}
		return config, nil
	}
	_, err := toml.DecodeFile(configFile, config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	if config.Themes == nil {
		config.Themes = make(map[string]models.Theme)
	}

	if config.Language == "" {
		config.Language = normalizeAppLanguage(GetSystemLocaleWindows())
	} else {
		config.Language = normalizeAppLanguage(config.Language)
	}

	config.ensureBuiltinThemes()
	if config.Radar.ThemeName == "" {
		config.Radar.ThemeName = models.DefaultTheme().Name
	}

	config.ensureRadarAnalysisDefaults()
	config.ensureLogsDefaults()

	return config, nil
}
