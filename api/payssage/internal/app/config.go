package app

import (
	"lib/database"
	"lib/errx"

	"github.com/caarlos0/env/v11"
	"github.com/google/uuid"
)

type Config struct {
	// Server
	Port      string `env:"PORT" envDefault:"8080"`
	AppName   string `env:"APP_NAME,required"`
	AppUrl    string `env:"APP_URL,required"`
	DebugMode bool   `env:"DEBUG_MODE"`

	// Migration
	MigrationPath string `env:"MIGRATION_PATH,required" envDefault:"./internal/database/migrations"`

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

	MpClientID        string `env:"MP_CLIENT_ID,required"`
	MpClientSecret    string `env:"MP_CLIENT_SECRET,required"`
	MpAccessToken     string `env:"MP_ACCESS_TOKEN,required"`
	MpRedirectURI     string `env:"MP_REDIRECT_URI,required"`
	MpWebhookSecret   string `env:"MP_WEBHOOK_SECRET,required"`
	MpTestAccessToken string `env:"MP_TEST_ACCESS_TOKEN,required"`
	MpTestPublicKey   string `env:"MP_TEST_PUBLIC_KEY,required"`

	// CORS
	CorsAllowedOrigins string `env:"CORS_ALLOWED_ORIGINS,required"`

	// Profiling
	ProfilePort string `env:"PROFILE_PORT"`

	// Feature flags
	DisableRateLimit bool `env:"DISABLE_RATE_LIMIT"`
}

func (cfg Config) ToDBConfig() database.Config {
	return database.Config{
		Host:          cfg.PostgresHost,
		Port:          cfg.PostgresPort,
		DB:            cfg.PostgresDB,
		User:          cfg.PostgresUser,
		Password:      cfg.PostgresPassword,
		SSLMode:       "disable",
		RootUser:      cfg.RootPostgresUser,
		RootPassword:  cfg.RootPostgresPassword,
		RootDB:        cfg.RootPostgresDB,
		RootHost:      "postgres",
		RootPort:      "5432",
		MigrationPath: cfg.MigrationPath,
	}
}

func LoadConfig() Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		errx.Exit(err, "failed to load config")
	}
	return cfg
}
