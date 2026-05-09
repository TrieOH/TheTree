package app

import (
	"github.com/caarlos0/env/v11"
	"github.com/google/uuid"
)

type Config struct {
	IdxURL             string    `env:"IDENTITY_X_URL,required"`
	IdxAPIKey          string    `env:"IDENTITY_X_API_KEY,required"`
	IdxProjectID       uuid.UUID `env:"IDENTITY_X_PROJECT_ID,required"`
	SpiceDBAddr        string    `env:"SPICEDB_ADDR,required"`
	SpiceDBToken       string    `env:"SPICEDB_TOKEN,required"`
	DatabaseURL        string    `env:"DATABASE_URL,required"`
	Port               string    `env:"PORT,required"`
	RedisAddr          string    `env:"REDIS_ADDR,required"`
	RedisPassword      string    `env:"REDIS_PASSWORD,required"`
	RedisDB            int       `env:"REDIS_DB,required"`
	AppName            string    `env:"APP_NAME,required"`
	CorsAllowedOrigins string    `env:"CORS_ALLOWED_ORIGINS,required"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	return cfg, env.Parse(&cfg)
}
