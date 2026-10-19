package config

import (
	"errors"
	"strconv"
	"time"

	libconfig "lib/config"
	"lib/database"
	"lib/email"
	"lib/errx"
	"lib/httpserver"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	libconfig.Server
	libconfig.CORS

	Postgres libconfig.Postgres `envPrefix:"IDX_"`

	// Server extras
	AppURL      string `env:"APP_URL,required"`
	ProfilePort string `env:"PROFILE_PORT"     envDefault:"6060"`

	// Email over SMTP (mailpit locally, any SMTP relay when deployed).
	EmailFrom    string `env:"EMAIL_FROM,required"`
	SMTPHost     string `env:"SMTP_HOST,required"`
	SMTPPort     string `env:"SMTP_PORT,required"`
	SMTPUser     string `env:"SMTP_USERNAME"`
	SMTPPass     string `env:"SMTP_PASSWORD"`
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
	return libconfig.DBConnectionConfig(cfg.Postgres)
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
		From:     cfg.EmailFrom,
		TLS:      cfg.SMTPTLS,
	}
}

func (cfg *Config) ToRateLimit() httpserver.RateLimit {
	return httpserver.RateLimit{
		Disabled: cfg.DisableRateLimit,
		RPS:      cfg.RateLimitRPS,
		Burst:    cfg.RateLimitBurst,
	}
}

// Validate checks the cross-field rules env tags cannot express.
func (cfg *Config) Validate() error {
	_, err := strconv.Atoi(cfg.SMTPPort)
	if err != nil {
		return errors.New("SMTP_PORT must be a number")
	}
	return nil
}

// LoadConfig parses the environment and validates it. Secrets arrive as
// plain environment variables (compose env_file locally, Worker secrets
// passed into the container on Cloudflare).
func LoadConfig() Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		errx.Exit(err, "error loading config")
	}
	err = cfg.Validate()
	if err != nil {
		errx.Exit(err, "invalid config")
	}
	return cfg
}
