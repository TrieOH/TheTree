package application

import (
	"time"
	"univents/internal/adapters/memory/imc"
	"univents/internal/adapters/memory/redis"
	"univents/internal/adapters/persistence"
	"univents/internal/infrastructure"

	"github.com/spf13/viper"
)

type Application struct {
}

func NewApplication(infra infrastructure.Infra) *Application {
	_ = persistence.NewRepositories(infra)

	cacheTTLStr := viper.GetString("KEYS_CACHE_TTL")
	cacheTTL, err := time.ParseDuration(cacheTTLStr)
	if err != nil {
		cacheTTL = time.Hour
	}

	_ = imc.NewInMemoryCache(100, cacheTTL)
	_ = imc.NewInMemoryCache(1000, cacheTTL)
	_ = redis.NewRedisCache(infra.Redis)

	return &Application{}
}
