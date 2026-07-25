package telemetry

import (
	"os"
	"sync"

	"lib/errx"

	"go.uber.org/zap"
)

var (
	logger      *zap.Logger
	debugLogger *zap.Logger
	once        sync.Once
)

func Init() {
	once.Do(func() {
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			errx.Exit(err, "error starting production logger")
		}
		debugLogger, err = zap.NewDevelopment()
		if err != nil {
			errx.Exit(err, "error starting development logger")
		}
	})
}

func Log() *zap.Logger {
	if logger == nil {
		Init()
	}
	if os.Getenv("TEST_MODE") == "TRUE" {
		return debugLogger
	}
	return logger
}
