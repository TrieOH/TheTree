package app

import (
	"time"

	"lib/database"
	"lib/errx"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	// Server
	Port        string `env:"PORT" envDefault:"8080"`
	ProfilePort string `env:"PROFILE_PORT" envDefault:"6060"`
	AppName     string `env:"APP_NAME,required"`
	AppUrl      string `env:"APP_URL,required"`
	DebugMode   bool   `env:"DEBUG_MODE"`

	// Postgres (own DB)
	PostgresHost     string `env:"IDX_POSTGRES_HOST,required"`
	PostgresPort     string `env:"IDX_POSTGRES_PORT" envDefault:"5432"`
	PostgresDB       string `env:"IDX_POSTGRES_DB,required"`
	PostgresUser     string `env:"IDX_POSTGRES_USER,required"`
	PostgresPassword string `env:"IDX_POSTGRES_PASSWORD,required"`

	// Migration
	MigrationPath string `env:"MIGRATION_PATH,required" envDefault:"./internal/database/migrations"`

	// Postgres (root — from .tree.env)
	RootPostgresUser     string `env:"POSTGRES_USER,required"`
	RootPostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	RootPostgresDB       string `env:"POSTGRES_DB" envDefault:"postgres"`

	// SMTP
	SmtpHost     string `env:"SMTP_HOST,required"`
	SmtpPort     string `env:"SMTP_PORT,required"`
	SmtpUser     string `env:"SMTP_USERNAME"`
	SmtpPass     string `env:"SMTP_PASSWORD"`
	SmtpFrom     string `env:"SMTP_FROM,required"`
	SmtpTls      bool   `env:"SMTP_TLS"`
	SmtpStartTls bool   `env:"SMTP_STARTTLS"`

	// Auth / crypto
	Issuer                string        `env:"ISSUER,required"`
	EncryptionKey         string        `env:"ENCRYPTION_KEY,required"`
	KeyLifetime           time.Duration `env:"IDENTITY_X_KEY_LIFETIME,required"`
	RotateKeysJobDuration time.Duration `env:"ROTATE_KEYS_JOB_DURATION,required"`

	// Tokens
	AccessTokenLifetime  time.Duration `env:"ACCESS_TOKEN_LIFETIME,required"`
	RefreshTokenLifetime time.Duration `env:"REFRESH_TOKEN_LIFETIME,required"`

	// CORS
	CorsAllowedOrigins string `env:"CORS_ALLOWED_ORIGINS,required"`
	CorsAllowedHeaders string `env:"CORS_ALLOWED_HEADERS,required"`

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
		errx.Exit(err, "error loading config")
	}
	return cfg
}
