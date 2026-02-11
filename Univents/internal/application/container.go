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
	repos := persistence.NewRepositories(infra)

	cacheTTLStr := viper.GetString("KEYS_CACHE_TTL")
	cacheTTL, err := time.ParseDuration(cacheTTLStr)
	if err != nil {
		cacheTTL = time.Hour
	}

	privateCache := imc.NewInMemoryCache(100, cacheTTL)
	publicCache := imc.NewInMemoryCache(1000, cacheTTL)

	sharedCache := redis.NewRedisCache(infra.Redis)

	return &Application{}
}
