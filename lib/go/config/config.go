// Package config holds the env blocks every Backend shares: the server
// settings, its own Postgres database, the IdentityX client, and CORS. Embed
// the blocks into the backend's Config; the own-DB block takes its service
// env prefix from the envPrefix tag, e.g. INFORMD_POSTGRES_HOST.
// Genuinely per-backend settings (feature fields, SMTP, providers) stay in
// the backend's config.go.
package config

import (
	"fmt"

	"lib/database"

	idx "sdk/identityx"

	"github.com/google/uuid"
)

// Server is the process-level block, unprefixed in every backend: the
// Runtime block, PORT, APP_NAME, DEBUG_MODE, and the rate limit knobs.
type Server struct {
	Runtime

	Port             string  `env:"PORT"               envDefault:"8080"`
	AppName          string  `env:"APP_NAME,required"`
	DebugMode        bool    `env:"DEBUG_MODE"`
	DisableRateLimit bool    `env:"DISABLE_RATE_LIMIT"`
	RateLimitRPS     float64 `env:"RATE_LIMIT_RPS"     envDefault:"400"`
	RateLimitBurst   int     `env:"RATE_LIMIT_BURST"   envDefault:"20"`
}

// RuntimeMode is how the Backend process is run.
//
//   - hosted: a long-lived HTTP server (compose, a VM). Background work runs
//     in-process next to the server.
//   - managed: a request-scoped function (AWS Lambda behind the Lambda Web
//     Adapter). The process may be frozen between requests, so background
//     work is handed to an external queue and scheduler instead.
type RuntimeMode string

const (
	HostedMode  RuntimeMode = "hosted"
	ManagedMode RuntimeMode = "managed"
)

// RuntimeRole is what a managed function does. Every role serves the same
// image; the role only decides what is mounted on top of the API.
//
//   - api: serves the HTTP API.
//   - worker: additionally serves the background event endpoint (queue
//     batches and scheduled jobs). Never exposed through the API gateway.
type RuntimeRole string

const (
	APIRole    RuntimeRole = "api"
	WorkerRole RuntimeRole = "worker"
)

// Runtime is the RUNTIME_MODE / RUNTIME_ROLE pair, unprefixed.
type Runtime struct {
	Mode RuntimeMode `env:"RUNTIME_MODE" envDefault:"hosted"`
	Role RuntimeRole `env:"RUNTIME_ROLE" envDefault:"api"`
}

// Managed reports whether the process runs as a managed function.
func (r Runtime) Managed() bool { return r.Mode == ManagedMode }

// Worker reports whether the process serves the background event endpoint.
func (r Runtime) Worker() bool { return r.Managed() && r.Role == WorkerRole }

// Validate rejects unknown modes and roles so a typo fails boot instead of
// silently falling back to a different process shape.
func (r Runtime) Validate() error {
	switch r.Mode {
	case HostedMode, ManagedMode:
	default:
		return fmt.Errorf("RUNTIME_MODE must be %q or %q, got %q", HostedMode, ManagedMode, r.Mode)
	}
	switch r.Role {
	case APIRole, WorkerRole:
	default:
		return fmt.Errorf("RUNTIME_ROLE must be %q or %q, got %q", APIRole, WorkerRole, r.Role)
	}
	return nil
}

// Postgres is the backend's own database, configured by POSTGRES_* in the
// service's environment.
type Postgres struct {
	Host           string `env:"POSTGRES_HOST,required"`
	Port           string `env:"POSTGRES_PORT"              envDefault:"5432"`
	DB             string `env:"POSTGRES_DB,required"`
	User           string `env:"POSTGRES_USER,required"`
	Password       string `env:"POSTGRES_PASSWORD,required"`
	SSLMode        string `env:"POSTGRES_SSLMODE"           envDefault:"require"`
	ChannelBinding string `env:"POSTGRES_CHANNEL_BINDING"   envDefault:"require"`
	// MaxConns caps the pool. Zero keeps pgx's default (max(4, NumCPU)).
	// Behind a pooler (Neon's pooled endpoint, pgbouncer) and in managed
	// functions keep it small: every function instance holds its own pool.
	MaxConns int32 `env:"POSTGRES_MAX_CONNS"`
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

// DBConfig assembles the service's direct database connection settings.
func DBConfig(p Postgres, migrationPath string) database.Config {
	cfg := DBConnectionConfig(p)
	cfg.MigrationPath = migrationPath
	return cfg
}

// DBConnectionConfig assembles database connection settings for a service
// whose schema lifecycle is managed outside the API process.
func DBConnectionConfig(p Postgres) database.Config {
	return database.Config{
		Host:           p.Host,
		Port:           p.Port,
		DB:             p.DB,
		User:           p.User,
		Password:       p.Password,
		SSLMode:        p.SSLMode,
		ChannelBinding: p.ChannelBinding,
		MaxConns:       p.MaxConns,
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
