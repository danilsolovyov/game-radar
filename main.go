package main

import (
	"context"
	"embed"
	"game-radar/app"
	"log/slog"
	"os"
	"sync"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

const (
	AppName = "Game Radar"
)

var appVersion = "dev"

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/windows/icon.ico
var appIcon []byte

func main() {
	slog.Info("main begin")
	config, err := app.LoadConfig()
	if err != nil {
		slog.Error("config load failed", "error", err)
		os.Exit(1)
	}
	slog.Info("config loaded")

	logs, err := setupLogging(config, AppName)
	if err != nil {
		slog.Error("logging setup failed", "error", err)
		os.Exit(1)
	}
	slog.Info("logging initialized")

	var closeLogsOnce sync.Once
	closeLogs := func() {
		closeLogsOnce.Do(func() {
			if closeErr := logs.closeFn(); closeErr != nil {
				slog.Error("logging close failed", "error", closeErr)
			}
		})
	}

	a := app.NewApp(AppName, appVersion, config, logs.backendLogger)
	slog.Info("app initialized", "version", appVersion)
	go systray.Run(func() { a.InitTray(appIcon) }, func() { a.ShutdownTray() })
	slog.Info("systray goroutine started")

	// Create application with options
	slog.Info("wails run start")
	err = wails.Run(&options.App{
		Title:  AppName,
		Width:  app.DefaultWindowWidth,
		Height: app.DefaultWindowHeight,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 0},
		AlwaysOnTop:      false,
		Frameless:        false,
		DisableResize:    false,
		OnStartup:        a.Startup,
		OnShutdown: func(_ context.Context) {
			slog.Info("onshutdown begin")
			a.StopRadar()
			a.ShutdownTray()
			systray.Quit()
			closeLogs()
			slog.Info("onshutdown complete")
		},
		Logger:             logs.wailsLogger,
		LogLevel:           logs.wailsLevel,
		LogLevelProduction: logs.wailsProd,
		Bind: []any{
			a,
		},
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               true, // enables window translucency in overlay mode
			DisableFramelessWindowDecorations: true,
			ResizeDebounceMS:                  5,
			Theme:                             windows.SystemDefault,
		},
	})

	if err != nil {
		slog.Error("wails run failed", "error", err)
		closeLogs()
		os.Exit(1)
	}
	slog.Info("wails run completed")
}
