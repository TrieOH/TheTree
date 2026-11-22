// Package config holds the env blocks the Backends share: the server
// settings, its own Postgres database, the IdentityX client, CORS, SMTP and
// background jobs. Embed the blocks into the backend's Config; the own-DB
// block takes its service env prefix from the envPrefix tag, e.g.
// INFORMD_POSTGRES_HOST. Genuinely per-backend settings (feature fields,
// providers) stay in the backend's config.go, and cmd/main reads the whole
// Config with Load.
package config

import (
	"uuid"

	"lib/database"
	"lib/email"
	"lib/errx"

	"github.com/caarlos0/env/v11"

	idx "sdk/identityx"
)

// Server is the process-level block, unprefixed in every backend: PORT,
// APP_NAME, API_URL, DEBUG_MODE, PROFILE_PORT, and the rate limit knobs.
type Server struct {
	Port    string `env:"PORT"              envDefault:"8080"`
	AppName string `env:"APP_NAME,required"`
	// APIURL is the backend's own public base URL; problem type URIs are
	// {APIURL}/problems/{slug}.
	APIURL           string  `env:"API_URL,required"`
	DebugMode        bool    `env:"DEBUG_MODE"`
	DisableRateLimit bool    `env:"DISABLE_RATE_LIMIT"`
	RateLimitRPS     float64 `env:"RATE_LIMIT_RPS"     envDefault:"400"`
	RateLimitBurst   int     `env:"RATE_LIMIT_BURST"   envDefault:"20"`
	// ProfilePort serves pprof when set; empty leaves it off.
	ProfilePort string `env:"PROFILE_PORT"`
	// ClientIPSecret is shared with our front Workers. A request carrying it
	// in X-Client-IP-Secret is trusted to name its user's IP in X-Client-IP,
	// so a front's server-side calls are rate limited per user instead of
	// all sharing the Worker's egress address. Empty: never trusted.
	ClientIPSecret string `env:"CLIENT_IP_SECRET"`
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
	// Behind a pooler (Neon's pooled endpoint, pgbouncer) keep it small:
	// every container instance holds its own pool.
	MaxConns int32 `env:"POSTGRES_MAX_CONNS"`
}

// IdentityX is the IdentityX client block (IDENTITY_X_URL / _API_KEY /
// _PROJECT_ID), unprefixed.
type IdentityX struct {
	URL       string    `env:"IDENTITY_X_URL,required"`
	APIKey    string    `env:"IDENTITY_X_API_KEY,required"`
	ProjectID uuid.UUID `env:"IDENTITY_X_PROJECT_ID,required"`
}

// CORS is the CORS pair, unprefixed. Both are comma-separated.
type CORS struct {
	AllowedOrigins string `env:"CORS_ALLOWED_ORIGINS,required"`
	AllowedHeaders string `env:"CORS_ALLOWED_HEADERS,required"`
}

// Jobs is the background jobs block, unprefixed. URL empty runs the jobs
// in-process (local compose); set, they're POSTed to that queue producer,
// which delivers them back to /_jobs (on Cloudflare the Worker sets it).
type Jobs struct {
	URL string `env:"JOBS_URL"`
	// Secret, when set, is required in the Jobs-Secret header of every
	// delivery to /_jobs. The API Worker sends it; nothing else can run a job
	// even if the path were ever reachable from outside.
	Secret string `env:"JOBS_SECRET"`
}

// SMTP is the outgoing mail block, unprefixed, for the backends that send
// email: mailpit locally, any SMTP relay when deployed.
type SMTP struct {
	Host     string `env:"SMTP_HOST,required"`
	Port     int    `env:"SMTP_PORT"          envDefault:"587"`
	Username string `env:"SMTP_USERNAME"`
	Password string `env:"SMTP_PASSWORD"`
	From     string `env:"SMTP_FROM,required"`
	TLS      bool   `env:"SMTP_TLS"           envDefault:"true"`
	// MonthlyLimit caps the recipients this service mails per calendar
	// month; 0 means no cap (local dev). The services' limits together stay
	// under Cloudflare's included 3,000 a month.
	MonthlyLimit int `env:"SMTP_MONTHLY_LIMIT"`
}

// EmailConfig assembles the SMTP client config from the shared block.
func EmailConfig(s SMTP) email.Config {
	return email.Config{
		Host:     s.Host,
		Port:     s.Port,
		Username: s.Username,
		Password: s.Password,
		From:     s.From,
		TLS:      s.TLS,
	}
}

// DBConnectionConfig assembles the service's database connection settings;
// the schema is migrated outside the API process.
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

// Load parses a backend's Config from the environment and ends the process
// when that fails: it belongs in cmd/main. Secrets arrive as plain
// environment variables (compose env_file locally, Worker secrets passed
// into the container on Cloudflare).
func Load[T any]() T {
	var cfg T
	err := env.Parse(&cfg)
	errx.Exit(err, "error loading config")
	return cfg
}
