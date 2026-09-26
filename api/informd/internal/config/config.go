package config

import (
	libconfig "lib/config"
	"lib/database"
	"lib/errx"

	idx "sdk/identityx"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	libconfig.Server
	libconfig.RootPostgres
	libconfig.IdentityX
	libconfig.CORS

	// Own database, prefixed INFORMD_.
	Postgres libconfig.Postgres `envPrefix:"INFORMD_"`

	// Migration
	MigrationPath string `env:"MIGRATION_PATH,required" envDefault:"./db/migrations"`

	// Feature fields
	AppURL      string `env:"APP_URL,required"`
	ProfilePort string `env:"PROFILE_PORT"`
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
