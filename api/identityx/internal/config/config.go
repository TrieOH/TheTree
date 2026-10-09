package config

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	libconfig "lib/config"
	"lib/database"
	"lib/email"
	"lib/errx"
	"lib/httpserver"
	"lib/secrets/provider"

	"github.com/caarlos0/env/v11"
)

// EmailTransport selects the email.Sender adapter.
type EmailTransport string

const (
	SMTPTransport EmailTransport = "smtp"
	SESTransport  EmailTransport = "ses"
)

type Config struct {
	libconfig.Server
	libconfig.CORS

	Postgres libconfig.Postgres `envPrefix:"IDX_"`

	// Server extras
	AppURL      string `env:"APP_URL,required"`
	ProfilePort string `env:"PROFILE_PORT"     envDefault:"6060"`

	// Email: the transport picks the adapter (smtp, or ses through the
	// standard AWS environment); EMAIL_FROM is the sender either way.
	EmailTransport EmailTransport `env:"EMAIL_TRANSPORT"     envDefault:"smtp"`
	EmailFrom      string         `env:"EMAIL_FROM,required"`

	// SMTP (EMAIL_TRANSPORT=smtp)
	SMTPHost     string `env:"SMTP_HOST"`
	SMTPPort     string `env:"SMTP_PORT"`
	SMTPUser     string `env:"SMTP_USERNAME"`
	SMTPPass     string `env:"SMTP_PASSWORD"`
	SMTPTLS      bool   `env:"SMTP_TLS"`
	SMTPStartTLS bool   `env:"SMTP_STARTTLS"`

	// Background jobs. Hosted mode runs them in-process; managed mode sends
	// them to this queue, drained by the worker role.
	JobsQueueURL string `env:"JOBS_QUEUE_URL"`

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
	errs := []error{cfg.Runtime.Validate()}
	switch cfg.EmailTransport {
	case SMTPTransport:
		if cfg.SMTPHost == "" || cfg.SMTPPort == "" {
			errs = append(errs, errors.New("EMAIL_TRANSPORT=smtp requires SMTP_HOST and SMTP_PORT"))
		}
	case SESTransport:
	default:
		errs = append(errs, fmt.Errorf("EMAIL_TRANSPORT must be %q or %q, got %q", SMTPTransport, SESTransport, cfg.EmailTransport))
	}
	if cfg.Managed() && cfg.JobsQueueURL == "" {
		errs = append(errs, errors.New("RUNTIME_MODE=managed requires JOBS_QUEUE_URL"))
	}
	return errors.Join(errs...)
}

// LoadConfig fills the environment from the configured secret store
// (SECRETS_PROVIDER; explicit environment wins), then parses and validates.
func LoadConfig() Config {
	err := provider.Load(context.Background())
	if err != nil {
		errx.Exit(err, "error loading secrets")
	}
	var cfg Config
	err = env.Parse(&cfg)
	if err != nil {
		errx.Exit(err, "error loading config")
	}
	err = cfg.Validate()
	if err != nil {
		errx.Exit(err, "invalid config")
	}
	return cfg
}
