package app

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	// Server
	Port      string `env:"PORT" envDefault:"8080"`
	AppName   string `env:"APP_NAME,required"`
	AppUrl    string `env:"APP_URL,required"`
	DebugMode bool   `env:"DEBUG_MODE"`

	// Postgres (own DB)
	PostgresHost     string `env:"IDX_POSTGRES_HOST,required"`
	PostgresPort     string `env:"IDX_POSTGRES_PORT" envDefault:"5432"`
	PostgresDB       string `env:"IDX_POSTGRES_DB,required"`
	PostgresUser     string `env:"IDX_POSTGRES_USER,required"`
	PostgresPassword string `env:"IDX_POSTGRES_PASSWORD,required"`

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

	// CORS
	CorsAllowedOrigins string `env:"CORS_ALLOWED_ORIGINS,required"`

	// Feature flags
	DisableRateLimit bool `env:"DISABLE_RATE_LIMIT"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	return cfg, env.Parse(&cfg)
}
