package logs

import (
	"go.uber.org/zap"
	"github.com/spf13/viper"
  "log"
	"sync"
)

var (
	logger *zap.Logger
	once   sync.Once
)

func Init() {
	once.Do(func() {
		var err error
		if viper.GetString("DEV_MODE") == "true" {
			logger, err = zap.NewDevelopment()
		} else {
			logger, err = zap.NewProduction()
		}
		if err != nil {
			log.Fatalf(err.Error())
		}
	})
}

func L() *zap.Logger {
	if logger == nil {
		Init()
	}
	return logger
}
