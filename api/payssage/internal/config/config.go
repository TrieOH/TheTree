package config

import (
	libconfig "lib/config"
	"lib/database"
	"lib/errx"
	"lib/httpserver"

	idx "sdk/identityx"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	MercadoPagoConfig

	libconfig.Server
	libconfig.IdentityX
	libconfig.CORS

	Postgres libconfig.Postgres `envPrefix:"PAYSSAGE_"`

	// Profiling
	ProfilePort string `env:"PROFILE_PORT"`

	// Background jobs. Empty: run them in-process (local compose). Set: POST
	// them to this queue producer, which delivers them back to /events; the
	// deploy's scheduler fires the periodic ones there.
	JobsURL string `env:"JOBS_URL"`
}

type MercadoPagoConfig struct {
	MpClientID        string `env:"MP_CLIENT_ID,required"`
	MpClientSecret    string `env:"MP_CLIENT_SECRET,required"`
	MpAccessToken     string `env:"MP_ACCESS_TOKEN,required"`
	MpRedirectURI     string `env:"MP_REDIRECT_URI,required"`
	MpWebhookSecret   string `env:"MP_WEBHOOK_SECRET,required"`
	MpTestAccessToken string `env:"MP_TEST_ACCESS_TOKEN,required"`
	MpTestPublicKey   string `env:"MP_TEST_PUBLIC_KEY,required"`
}

func (cfg Config) ToIdentityXConfig() idx.Config {
	return libconfig.IdentityXConfig(cfg.IdentityX, true)
}

// ToDBConfig is the connection only: the schema is migrated outside the
// API process (`just payssage-goose`).
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

func LoadConfig() Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		errx.Exit(err, "failed to load config")
	}
	return cfg
}
