package config

import (
	libconfig "lib/config"
	"lib/database"
	"lib/email"
	"lib/errx"
	"lib/httpserver"

	idx "sdk/identityx"

	"github.com/caarlos0/env/v11"
	"github.com/google/uuid"
)

type Config struct {
	libconfig.Server
	libconfig.IdentityX
	libconfig.CORS

	Postgres libconfig.Postgres `envPrefix:"UNIVENTS_"`

	// Profiling
	ProfilePort string `env:"PROFILE_PORT" envDefault:"6060"`

	// Feature fields
	AppURL string `env:"APP_URL,required"`

	// When true, every JWT-authenticated route requires the actor's email to
	// be verified (identityx subject.verified_at set). Unverified actors get
	// 403. Public routes (webhooks, health, ws) are unaffected.
	RequireVerifiedEmail bool `env:"REQUIRE_VERIFIED_EMAIL" envDefault:"false"`

	// Payssage
	PayssageURL           string    `env:"PAYSSAGE_URL,required"`
	PayssageAPIKey        string    `env:"PAYSSAGE_API_KEY,required"`
	PayssageWalletID      uuid.UUID `env:"PAYSSAGE_WALLET_ID,required"` // one platform wallet shared by every event (D6)
	PayssageWebhookSecret string    `env:"PAYSSAGE_WEBHOOK_SECRET,required"`

	// Object Storage (RustFS)
	ObjStorageEndpoint  string `env:"OBJECT_STORAGE_ENDPOINT,required"`
	ObjStorageAccessKey string `env:"OBJECT_STORAGE_ACCESS_KEY,required"`
	ObjStorageSecretKey string `env:"OBJECT_STORAGE_SECRET_KEY,required"`
	ObjStorageUseSSL    bool   `env:"OBJECT_STORAGE_USE_SSL"             envDefault:"true"`
	ObjStorageRegion    string `env:"OBJECT_STORAGE_REGION"              envDefault:"us-east-1"`

	// SMTP / Email
	SMTPHost     string `env:"SMTP_HOST,required"`
	SMTPPort     int    `env:"SMTP_PORT"          envDefault:"587"`
	SMTPUsername string `env:"SMTP_USERNAME"`
	SMTPPassword string `env:"SMTP_PASSWORD"`
	SMTPFrom     string `env:"SMTP_FROM,required"`
	SMTPTLS      bool   `env:"SMTPTLS"            envDefault:"true"`

	// Webhook HMAC
	HmacSecret string `env:"HMAC_SECRET,required"`

	// Background jobs. Empty: run them in-process (local compose). Set: POST
	// them to this queue producer, which delivers them back to /events.
	JobsURL string `env:"JOBS_URL"`
}

func (cfg Config) ToIdentityXConfig() idx.Config {
	return libconfig.IdentityXConfig(cfg.IdentityX, cfg.DebugMode)
}

// ToDBConfig is the connection only: the schema is migrated outside the
// API process (`just univents-goose`).
func (cfg Config) ToDBConfig() database.Config {
	return libconfig.DBConnectionConfig(cfg.Postgres)
}

func (cfg Config) ToRateLimit() httpserver.RateLimit {
	return httpserver.RateLimit{
		Disabled: cfg.DisableRateLimit,
		RPS:      cfg.RateLimitRPS,
		Burst:    cfg.RateLimitBurst,
	}
}

func (cfg Config) ToEmailConfig() email.Config {
	return email.Config{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Username: cfg.SMTPUsername,
		Password: cfg.SMTPPassword,
		From:     cfg.SMTPFrom,
		TLS:      cfg.SMTPTLS,
	}
}

func Load() Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		errx.Exit(err, "error loading config")
	}
	return cfg
}
