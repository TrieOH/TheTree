package config

import (
	libconfig "lib/config"
	"lib/database"
	"lib/errx"

	idx "sdk/identityx"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	MercadoPagoConfig

	libconfig.Server
	libconfig.RootPostgres
	libconfig.IdentityX
	libconfig.CORS

	// Own database, prefixed PAYSSAGE_.
	Postgres libconfig.Postgres `envPrefix:"PAYSSAGE_"`

	// Migration
	MigrationPath string `env:"MIGRATION_PATH,required" envDefault:"./db/migrations"`

	// Profiling
	ProfilePort string `env:"PROFILE_PORT"`

	// Simple auth (e.g. River UI dashboard)
	SimpleAuthUser     string `env:"SIMPLE_AUTH_USER,required"`
	SimpleAuthPassword string `env:"SIMPLE_AUTH_PASS,required"`
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

func (cfg Config) ToDBConfig() database.Config {
	return libconfig.DBConfig(cfg.Postgres, cfg.RootPostgres, cfg.MigrationPath)
}

func LoadConfig() Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		errx.Exit(err, "failed to load config")
	}
	return cfg
}
