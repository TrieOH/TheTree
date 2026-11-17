package telemetry

import (
	"log/slog"
	"os"
	"sync"
)

var (
	logger *slog.Logger
	once   sync.Once
)

func Init() {
	once.Do(func() {
		logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{AddSource: true}))
	})
}

func Log() *slog.Logger {
	if logger == nil {
		Init()
	}
	return logger
}

// Err is the attr errors are logged under.
func Err(err error) slog.Attr {
	return slog.Any("error", err)
}
