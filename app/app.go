package app

import (
	"context"
	"game-radar/services/audiodriver"
	"game-radar/services/radar"
	"log/slog"
	"runtime/debug"
	"sync"

	"github.com/getlantern/systray"
)

// App struct.
type App struct {
	radarCancelFn context.CancelFunc
	ctx           context.Context
	driver        audiodriver.AudioDriver
	radar         *radar.Radar
	appName       string
	appVersion    string
	config        *Config
	logger        *slog.Logger
	mu            sync.Mutex
	overlayMode   bool
	overlayStatus chan bool

	trayTitle      string
	trayIcon       []byte
	trayToggleItem *systray.MenuItem
	trayQuitItem   *systray.MenuItem
	trayStopCh     chan struct{}
}

// NewApp creates a new App application struct.
func NewApp(appName, appVersion string, config *Config, logger *slog.Logger) *App {
	if logger == nil {
		logger = slog.Default()
	}
	wca := audiodriver.NewWCA(logger)

	return &App{
		driver: wca,
		radar: radar.NewRadarService(
			wca,
			config.Radar.IntensityFilter,
			config.Radar.Amplifier,
		),
		config:        config,
		appName:       appName,
		appVersion:    appVersion,
		logger:        logger,
		overlayStatus: make(chan bool, 8),
		trayTitle:     appName,
		trayStopCh:    make(chan struct{}),
	}
}

// so we can call the runtime methods.
func (a *App) Startup(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			a.logger.Error("startup panic recovered", "panic", r, "stack", string(debug.Stack()))
		}
	}()
	a.logger.Info("startup begin")
	a.ctx = ctx
	a.logger.Info("startup step: overlay status broadcaster")
	a.startOverlayStatusBroadcaster()
	a.logger.Info("startup step: set normal window")
	a.SetNormalWindow()
	a.logger.Info("startup step: emit language changed", "language", a.GetLanguage())
	a.emitLanguageChanged(a.GetLanguage())
	a.logger.Info("startup step: start radar async")
	a.startRadarAsync()
	a.logger.Info("startup complete")
}

func (a *App) GetAppName() string {
	return a.appName
}

func (a *App) GetVersion() string {
	return a.appVersion
}
