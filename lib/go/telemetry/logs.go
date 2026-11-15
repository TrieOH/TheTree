package telemetry

import (
	"log/slog"
	"os"
	"sync"
)

var (
	logger      *slog.Logger
	debugLogger *slog.Logger
	once        sync.Once
)

func Init() {
	once.Do(func() {
		logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{AddSource: true}))
		debugLogger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}))
	})
}

func Log() *slog.Logger {
	if logger == nil {
		Init()
	}
	if os.Getenv("TEST_MODE") == "TRUE" {
		return debugLogger
	}
	return logger
}

// Err is the attr errors are logged under.
func Err(err error) slog.Attr {
	return slog.Any("error", err)
}
