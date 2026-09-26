package config

import (
	libconfig "lib/config"
	"lib/database"
	"lib/email"
	"lib/errx"

	idx "sdk/identityx"

	"github.com/caarlos0/env/v11"
	"github.com/google/uuid"
)

type Config struct {
	libconfig.Server
	libconfig.RootPostgres
	libconfig.IdentityX
	libconfig.CORS

	// Own database, prefixed UNIVENTS_.
	Postgres libconfig.Postgres `envPrefix:"UNIVENTS_"`

	// Migration
	MigrationPath string `env:"MIGRATION_PATH,required" envDefault:"./db/migrations"`

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
}

func (cfg Config) ToIdentityXConfig() idx.Config {
	return libconfig.IdentityXConfig(cfg.IdentityX, cfg.DebugMode)
}

func (cfg Config) ToDBConfig() database.Config {
	return libconfig.DBConfig(cfg.Postgres, cfg.RootPostgres, cfg.MigrationPath)
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
