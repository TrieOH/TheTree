package telemetry

import (
	"lib/errx"
	"sync"

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
			errx.Must(err, "error starting production logger")
		}
		debugLogger, err = zap.NewDevelopment()
		if err != nil {
			errx.Must(err, "error starting development logger")
		}
	})
}

func Log() *zap.Logger {
	if logger == nil {
		Init()
	}
	return logger
}

func DLog() *zap.Logger {
	if debugLogger == nil {
		Init()
	}
	return debugLogger
}
