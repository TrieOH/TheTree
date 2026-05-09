package telemetry

import (
	"Informd/internal/shared/errx"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConfig struct {
	Level       string // "debug", "info", "warn", "error"
	Development bool
}

func NewLogger(cfg LogConfig) *zap.Logger {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		errx.Must(fmt.Errorf("invalid log level %q: %w", cfg.Level, err), "error setting log level on logger")
	}

	var zapCfg zap.Config
	if cfg.Development {
		zapCfg = zap.NewDevelopmentConfig()
	} else {
		zapCfg = zap.NewProductionConfig()
	}

	zapCfg.Level = zap.NewAtomicLevelAt(level)
	logger, err := zapCfg.Build()
	if err != nil {
		errx.Must(err, "error building logger")
	}
	return logger
}
