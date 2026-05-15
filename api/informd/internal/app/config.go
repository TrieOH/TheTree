package app

import (
	"github.com/caarlos0/env/v11"
	"github.com/google/uuid"
)

type Config struct {
	// Server
	Port      string `env:"PORT" envDefault:"8080"`
	AppName   string `env:"APP_NAME,required"`
	AppUrl    string `env:"APP_URL,required"`
	DebugMode bool   `env:"DEBUG_MODE"`

	// Postgres (own DB)
	PostgresHost     string `env:"INFORMD_POSTGRES_HOST,required"`
	PostgresPort     string `env:"INFORMD_POSTGRES_PORT" envDefault:"5432"`
	PostgresDB       string `env:"INFORMD_POSTGRES_DB,required"`
	PostgresUser     string `env:"INFORMD_POSTGRES_USER,required"`
	PostgresPassword string `env:"INFORMD_POSTGRES_PASSWORD,required"`

	// Postgres (root — from .tree.env)
	RootPostgresUser     string `env:"POSTGRES_USER,required"`
	RootPostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	RootPostgresDB       string `env:"POSTGRES_DB" envDefault:"postgres"`

	// Identity-X
	IdxURL       string    `env:"IDENTITY_X_URL,required"`
	IdxAPIKey    string    `env:"IDENTITY_X_API_KEY,required"`
	IdxProjectID uuid.UUID `env:"IDENTITY_X_PROJECT_ID,required"`

	// CORS
	CorsAllowedOrigins string `env:"CORS_ALLOWED_ORIGINS,required"`

	// Feature flags
	DisableRateLimit bool `env:"DISABLE_RATE_LIMIT"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	return cfg, env.Parse(&cfg)
}
