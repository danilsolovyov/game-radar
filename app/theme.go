package app

import "game-radar/models"

// GetThemes returns the list of available themes.
func (a *App) GetThemes() ([]models.Theme, error) {
	return a.config.getThemes(), nil
}

// SetTheme creates or updates a theme and makes it active.
func (a *App) SetTheme(theme models.Theme) error {
	return a.config.setTheme(theme)
}

// DeleteTheme removes a theme by name.
func (a *App) DeleteTheme(name string) error {
	return a.config.deleteTheme(name)
}

// GetCurrentThemeName returns the active theme name.
func (a *App) GetCurrentThemeName() string {
	return a.config.getCurrentThemeName()
}
