package main

import (
	"fmt"
	"game-radar/app"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	wailslogger "github.com/wailsapp/wails/v2/pkg/logger"
	"gopkg.in/natefinch/lumberjack.v2"
)

type logSetup struct {
	backendLogger *slog.Logger
	wailsLogger   wailslogger.Logger
	wailsLevel    wailslogger.LogLevel
	wailsProd     wailslogger.LogLevel
	closeFn       func() error
}

type rotatingWailsLogger struct {
	mu     sync.Mutex
	writer io.Writer
}

func (l *rotatingWailsLogger) Print(message string) {
	l.writeRaw(message)
}

func (l *rotatingWailsLogger) Trace(message string) {
	l.writeLevel("TRACE", message)
}

func (l *rotatingWailsLogger) Debug(message string) {
	l.writeLevel("DEBUG", message)
}

func (l *rotatingWailsLogger) Info(message string) {
	l.writeLevel("INFO", message)
}

func (l *rotatingWailsLogger) Warning(message string) {
	l.writeLevel("WARN", message)
}

func (l *rotatingWailsLogger) Error(message string) {
	l.writeLevel("ERROR", message)
}

func (l *rotatingWailsLogger) Fatal(message string) {
	l.writeLevel("FATAL", message)
	os.Exit(1)
}

func (l *rotatingWailsLogger) writeRaw(message string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	_, _ = fmt.Fprint(l.writer, message)
}

func (l *rotatingWailsLogger) writeLevel(level, message string) {
	timestamp := time.Now().Format(time.RFC3339)
	l.writeRaw(fmt.Sprintf("%s %s %s\n", timestamp, level, message))
}

func setupLogging(cfg *app.Config, appName string) (*logSetup, error) {
	logDir := resolveLogsDir(cfg.Logs.Dir, appName)
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, fmt.Errorf("create log directory: %w", err)
	}

	backendWriter := newRotatingWriter(logDir, cfg.Logs.Backend.File, cfg)
	wailsWriter := newRotatingWriter(logDir, cfg.Logs.Wails.File, cfg)
	backendSink := io.Writer(backendWriter)
	wailsSink := io.Writer(wailsWriter)
	if isWailsDevMode() {
		// In `wails dev` mode, duplicate output to the terminal.
		backendSink = io.MultiWriter(backendWriter, os.Stdout)
		wailsSink = io.MultiWriter(wailsWriter, os.Stdout)
	}

	backendLevel := parseSlogLevel(cfg.Logs.Backend.Level)
	slogHandler := slog.NewTextHandler(backendSink, &slog.HandlerOptions{
		Level:     backendLevel,
		AddSource: false,
	})
	backendLogger := slog.New(slogHandler)
	slog.SetDefault(backendLogger)

	wailsLogger := &rotatingWailsLogger{writer: wailsSink}

	return &logSetup{
		backendLogger: backendLogger,
		wailsLogger:   wailsLogger,
		wailsLevel:    parseWailsLevel(cfg.Logs.Wails.Level, wailslogger.INFO),
		wailsProd:     parseWailsLevel(cfg.Logs.Wails.LevelProduction, wailslogger.ERROR),
		closeFn: func() error {
			var closeErrs []string
			for _, c := range []io.Closer{backendWriter, wailsWriter} {
				if err := c.Close(); err != nil {
					closeErrs = append(closeErrs, err.Error())
				}
			}
			if len(closeErrs) > 0 {
				return fmt.Errorf("close logs: %s", strings.Join(closeErrs, "; "))
			}
			return nil
		},
	}, nil
}

func isWailsDevMode() bool {
	devServer := strings.ToLower(strings.TrimSpace(os.Getenv("devserver")))
	return devServer != "" && devServer != "0" && devServer != "false" && devServer != "off"
}

func newRotatingWriter(logDir, filename string, cfg *app.Config) *lumberjack.Logger {
	name := strings.TrimSpace(filename)
	if name == "" {
		name = "app.log"
	}
	return &lumberjack.Logger{
		Filename:   filepath.Join(logDir, name),
		MaxSize:    cfg.Logs.Rotation.MaxSizeMB,
		MaxBackups: cfg.Logs.Rotation.MaxBackups,
		MaxAge:     cfg.Logs.Rotation.MaxAgeDays,
		Compress:   cfg.Logs.Rotation.Compress,
		LocalTime:  true,
	}
}

func resolveLogsDir(configuredDir, appName string) string {
	dir := strings.TrimSpace(configuredDir)
	base := defaultLogsBaseDir(appName)

	if dir == "" || dir == app.DefaultLogsDir {
		return filepath.Join(base, app.DefaultLogsDir)
	}
	if filepath.IsAbs(dir) {
		return dir
	}
	return filepath.Join(base, dir)
}

func defaultLogsBaseDir(appName string) string {
	localAppData := strings.TrimSpace(os.Getenv("LOCALAPPDATA"))
	if localAppData != "" {
		return filepath.Join(localAppData, appName)
	}
	userConfigDir, err := os.UserConfigDir()
	if err == nil && userConfigDir != "" {
		return filepath.Join(userConfigDir, appName)
	}
	return "."
}

func parseSlogLevel(level string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func parseWailsLevel(level string, fallback wailslogger.LogLevel) wailslogger.LogLevel {
	parsed, err := wailslogger.StringToLogLevel(level)
	if err != nil {
		return fallback
	}
	return parsed
}
