package config

import (
	"strconv"
	"time"

	libconfig "lib/config"
	"lib/database"
	"lib/email"
	"lib/errx"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	libconfig.Server
	libconfig.RootPostgres
	libconfig.CORS

	// Own database, prefixed IDX_.
	Postgres libconfig.Postgres `envPrefix:"IDX_"`

	// Migration
	MigrationPath string `env:"MIGRATION_PATH,required" envDefault:"./db/migrations"`

	// Server extras
	AppURL      string `env:"APP_URL,required"`
	ProfilePort string `env:"PROFILE_PORT"     envDefault:"6060"`

	// SMTP
	SMTPHost     string `env:"SMTP_HOST,required"`
	SMTPPort     string `env:"SMTP_PORT,required"`
	SMTPUser     string `env:"SMTP_USERNAME"`
	SMTPPass     string `env:"SMTP_PASSWORD"`
	SMTPFrom     string `env:"SMTP_FROM,required"`
	SMTPTLS      bool   `env:"SMTP_TLS"`
	SMTPStartTLS bool   `env:"SMTP_STARTTLS"`

	// Auth / crypto
	Issuer                string        `env:"ISSUER,required"`
	EncryptionKey         string        `env:"ENCRYPTION_KEY,required"`
	HmacSecret            string        `env:"HMAC_SECRET,required"`
	KeyLifetime           time.Duration `env:"IDENTITY_X_KEY_LIFETIME,required"`
	RotateKeysJobDuration time.Duration `env:"ROTATE_KEYS_JOB_DURATION,required"`

	// Tokens
	AccessTokenLifetime  time.Duration `env:"ACCESS_TOKEN_LIFETIME,required"`
	RefreshTokenLifetime time.Duration `env:"REFRESH_TOKEN_LIFETIME,required"`

	// Action tokens (email verify / password reset links)
	EmailVerifyTokenTTL time.Duration `env:"EMAIL_VERIFY_TOKEN_TTL" envDefault:"10m"`
	EmailResetTokenTTL  time.Duration `env:"EMAIL_RESET_TOKEN_TTL"  envDefault:"10m"`
}

func (cfg *Config) ToDBConfig() database.Config {
	return libconfig.DBConfig(cfg.Postgres, cfg.RootPostgres, cfg.MigrationPath)
}

func (cfg *Config) ToEmailConfig() email.Config {
	port, err := strconv.Atoi(cfg.SMTPPort)
	if err != nil {
		port = 587
	}
	return email.Config{
		Host:     cfg.SMTPHost,
		Port:     port,
		Username: cfg.SMTPUser,
		Password: cfg.SMTPPass,
		From:     cfg.SMTPFrom,
		TLS:      cfg.SMTPTLS,
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
