// Package config holds the env blocks every Backend shares: the server
// settings, its own Postgres database, the root Postgres used for
// migrations, the IdentityX client, and CORS. Embed the blocks into the
// backend's Config; the own-DB block takes the service's env prefix via
// the envPrefix tag ("POSTGRES_HOST" + "INFORMD_" → INFORMD_POSTGRES_HOST).
// Genuinely per-backend settings (feature fields, SMTP, providers) stay in
// the backend's config.go.
package config

import (
	"lib/database"

	idx "sdk/identityx"

	"github.com/google/uuid"
)

// Server is the process-level block, unprefixed in every backend: PORT,
// APP_NAME, DEBUG_MODE, DISABLE_RATE_LIMIT.
type Server struct {
	Port             string `env:"PORT"               envDefault:"8080"`
	AppName          string `env:"APP_NAME,required"`
	DebugMode        bool   `env:"DEBUG_MODE"`
	DisableRateLimit bool   `env:"DISABLE_RATE_LIMIT"`
}

// Postgres is the backend's own database. Declare it with the service's
// env prefix: `Postgres Postgres \`envPrefix:"INFORMD_"\“.
type Postgres struct {
	Host     string `env:"POSTGRES_HOST,required"`
	Port     string `env:"POSTGRES_PORT"              envDefault:"5432"`
	DB       string `env:"POSTGRES_DB,required"`
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
}

// RootPostgres is the shared root database used for migrations
// (POSTGRES_USER / POSTGRES_PASSWORD / POSTGRES_DB), unprefixed.
type RootPostgres struct {
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	DB       string `env:"POSTGRES_DB"                envDefault:"postgres"`
}

// IdentityX is the IdentityX client block (IDENTITY_X_URL / _API_KEY /
// _PROJECT_ID), unprefixed.
type IdentityX struct {
	URL       string    `env:"IDENTITY_X_URL,required"`
	APIKey    string    `env:"IDENTITY_X_API_KEY,required"`
	ProjectID uuid.UUID `env:"IDENTITY_X_PROJECT_ID,required"`
}

// CORS is the CORS pair, unprefixed.
type CORS struct {
	AllowedOrigins string `env:"CORS_ALLOWED_ORIGINS,required"`
	AllowedHeaders string `env:"CORS_ALLOWED_HEADERS,required"`
}

// DBConfig assembles the lib/database config from the shared blocks. The
// connection-shaping policy (SSL mode, root host/port) lives here, once.
func DBConfig(p Postgres, r RootPostgres, migrationPath string) database.Config {
	return database.Config{
		Host:          p.Host,
		Port:          p.Port,
		DB:            p.DB,
		User:          p.User,
		Password:      p.Password,
		SSLMode:       "disable",
		RootUser:      r.User,
		RootPassword:  r.Password,
		RootDB:        r.DB,
		RootHost:      "postgres",
		RootPort:      "5432",
		MigrationPath: migrationPath,
	}
}

// IdentityXConfig assembles the SDK client config from the shared block.
func IdentityXConfig(i IdentityX, debug bool) idx.Config {
	return idx.Config{
		BaseURL:   i.URL,
		APIKey:    i.APIKey,
		ProjectID: i.ProjectID,
		Debug:     debug,
	}
}
