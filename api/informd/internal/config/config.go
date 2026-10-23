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
	libconfig.Server
	libconfig.IdentityX
	libconfig.CORS

	Postgres libconfig.Postgres `envPrefix:"INFORMD_"`

	// Feature fields
	AppURL      string `env:"APP_URL,required"`
	ProfilePort string `env:"PROFILE_PORT"`
}

func (cfg Config) ToIdentityXConfig() idx.Config {
	return libconfig.IdentityXConfig(cfg.IdentityX, true)
}

// ToDBConfig is the connection only: the schema is migrated outside the
// API process (`just informd-goose`).
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
