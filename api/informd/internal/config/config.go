package config

import (
	libconfig "lib/config"
	"lib/errx"

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

func LoadConfig() Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		errx.Exit(err, "failed to load config")
	}
	return cfg
}
