package app

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseURL           string        `env:"DATABASE_URL,required"`
	Port                  string        `env:"PORT,required"`
	DebugMode             bool          `env:"DEBUG_MODE"`
	AppName               string        `env:"APP_NAME,required"`
	SmtpHost              string        `env:"SMTP_HOST,required"`
	SmtpPort              string        `env:"SMTP_PORT,required"`
	SmtpUser              string        `env:"SMTP_USERNAME,required"`
	SmtpPass              string        `env:"SMTP_PASSWORD,required"`
	SmtpFrom              string        `env:"SMTP_FROM,required"`
	SmtpTls               bool          `env:"SMTP_TLS,required"`
	SmtpStartTls          bool          `env:"SMTP_STARTTLS,required"`
	CorsAllowedOrigins    string        `env:"CORS_ALLOWED_ORIGINS,required"`
	Issuer                string        `env:"ISSUER,required"`
	AppUrl                string        `env:"APP_URL,required"`
	DisableRateLimit      bool          `env:"DISABLE_RATE_LIMIT,required"`
	EncryptionKey         string        `env:"ENCRYPTION_KEY,required"`
	RedisAddress          string        `env:"REDIS_ADDR,required"`
	RedisPassword         string        `env:"REDIS_PASSWORD,required"`
	RedisDB               int           `env:"REDIS_DB,required"`
	KeyLifetime           time.Duration `env:"IDENTITY_X_KEY_LIFETIME,required"`
	RotateKeysJobDuration time.Duration `env:"ROTATE_KEYS_JOB_DURATION,required"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	return cfg, env.Parse(&cfg)
}
