package config

import (
	libconfig "lib/config"
)

type Config struct {
	libconfig.Server
	libconfig.IdentityX
	libconfig.CORS

	Postgres libconfig.Postgres `envPrefix:"INFORMD_"`

	// Feature fields
	AppURL string `env:"APP_URL,required"`
}
