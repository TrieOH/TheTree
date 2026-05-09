package app

import (
	"github.com/caarlos0/env/v11"
	"github.com/google/uuid"
)

type Config struct {
	IdxURL                string    `env:"IDENTITY_X_URL,required"`
	IdxAPIKey             string    `env:"IDENTITY_X_API_KEY,required"`
	IdxProjectID          uuid.UUID `env:"IDENTITY_X_PROJECT_ID,required"`
	PayssageURL           string    `env:"PAYSSAGE_URL,required"`
	PayssageProvider      string    `env:"PAYSSAGE_PROVIDER,required"`
	PayssageAPIKey        string    `env:"PAYSSAGE_API_KEY,required"`
	PayssageWebhookSecret string    `env:"PAYSSAGE_WEBHOOK_SECRET,required"`
	SpiceDBAddr           string    `env:"SPICEDB_ADDR,required"`
	SpiceDBToken          string    `env:"SPICEDB_TOKEN,required"`
	WsJwtSecret           string    `env:"WS_JWT_SECRET,required"`
	DatabaseURL           string    `env:"DATABASE_URL,required"`
	Port                  string    `env:"PORT,required"`
	RedisAddr             string    `env:"REDIS_ADDR,required"`
	RedisPassword         string    `env:"REDIS_PASSWORD,required"`
	RedisDB               int       `env:"REDIS_DB,required"`
	AppName               string    `env:"APP_NAME,required"`
	CorsAllowedOrigins    string    `env:"CORS_ALLOWED_ORIGINS,required"`
	ObjStorageURL         string    `env:"OBJECT_STORAGE_ENDPOINT,required"`
	ObjStorageAccessKey   string    `env:"OBJECT_STORAGE_ACCESS_KEY,required"`
	ObjStorageSecretKey   string    `env:"OBJECT_STORAGE_SECRET_KEY,required"`
	ObjStorageUseSSL      bool      `env:"OBJECT_STORAGE_USE_SSL,required"`
	DisableRateLimit      bool      `env:"DISABLE_RATE_LIMIT,required"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	return cfg, env.Parse(&cfg)
}
